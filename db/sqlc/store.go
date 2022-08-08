package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Store interface enables both the MockDB and our real DB can use this queries
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	EntryTx(ctx context.Context, arg EntryTxParams) (EntryTxResult, error)
}

//* Store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

//* TransferTxParams hold the all necessary input values for the transfer
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//* TransferTxResult holds the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type EntryTxParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type EntryTxResult struct {
	Entry   Entry   `json:"entry"`
	Account Account `json:"account"`
}

//* NewStore returns a new Store with the given DB
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

//* execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

//* TransferTx performs a money transfer from one account to the other.
//* It creates a transfer record, add account entries, and update account balances within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//! to avoid DB deadlock reorganized the order of the DB funcs

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

			if err != nil {
				return err
			}

		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

//* EntryTx performs a money entry for an account.
//* It checks the funds of the account, updates the account and creates a new entry
func (store *SQLStore) EntryTx(ctx context.Context, arg EntryTxParams) (EntryTxResult, error) {
	var result EntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		acc, err := q.GetAccount(ctx, arg.AccountID)

		if err != nil {
			return err
		}

		if arg.Amount < 0 && (-arg.Amount) > acc.Balance {
			err = errors.New("insufficient funds")
			return err
		}

		arg1 := AddAccountBalanceParams{
			Amount: arg.Amount,
			ID:     arg.AccountID,
		}

		result.Account, err = q.AddAccountBalance(ctx, arg1)

		if err != nil {
			return err
		}

		arg2 := CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    arg.Amount,
		}

		result.Entry, err = q.CreateEntry(ctx, arg2)

		if err != nil {
			return err
		}

		return err
	})

	return result, err
}

//* addMoney adds an amount to two given accounts
func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
