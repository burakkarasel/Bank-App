package db

import (
	"context"
	"testing"
	"time"

	"github.com/burakkarasel/Bank-App/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	testTransfer := createRandomTransfer(t, acc1, acc2)

	transfer, err := testQueries.GetTransfer(context.Background(), testTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, testTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, testTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, testTransfer.ID, transfer.ID)
	require.Equal(t, testTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, testTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	arg := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
	}
}
