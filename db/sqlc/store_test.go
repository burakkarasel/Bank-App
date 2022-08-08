package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTransferTx tests the TransferTx DB func
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	//* run n concurrent transfer transactions

	n := 5

	amount := int64(10)

	errs := make(chan error)

	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//! check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//! check transfer
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//! check from entry
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		require.NoError(t, err)

		//! check to entry
		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)

		require.NoError(t, err)

		//! check accounts

		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		//! check account's balance

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //! check if the diff equalt to transfer amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//! check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, acc1.Balance-amount*int64(n), updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+amount*int64(n), updatedAccount2.Balance)
}

// TestTransferTxDeadlock tests the TransferTx DB func for deadlocks
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	//* run n concurrent transfer transactions

	n := 10

	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 1 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	//! check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	//! check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, acc1.Balance, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance, updatedAccount2.Balance)
}

// TestEntryTx tests EntryTx transition func
func TestEntryTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)

	//* run n concurrent transfer transactions

	n := 5

	amount := int64(10)

	errs := make(chan error)

	results := make(chan EntryTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.EntryTx(context.Background(), EntryTxParams{
				Amount:    amount,
				AccountID: acc1.ID,
			})

			errs <- err
			results <- result
		}()
	}

	//! check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//! check entry
		entry := result.Entry

		require.NotEmpty(t, entry)
		require.Equal(t, acc1.ID, entry.AccountID)
		require.Equal(t, amount, entry.Amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry.ID)

		require.NoError(t, err)

		//! check account

		account := result.Account

		require.NotEmpty(t, account)
		require.Equal(t, acc1.ID, account.ID)
		require.Equal(t, acc1.Owner, account.Owner)
		require.Equal(t, acc1.Currency, account.Currency)
		require.WithinDuration(t, acc1.CreatedAt, account.CreatedAt, time.Second)
	}

	//! check the final updated balances
	updatedAccount, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, acc1.Balance+amount*int64(n), updatedAccount.Balance)
}
