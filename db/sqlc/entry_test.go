package db

import (
	"context"
	"testing"
	"time"

	"github.com/burakkarasel/Bank-App/util"
	"github.com/stretchr/testify/require"
)

// createRandomEntry creates a random entry
func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NoError(t, err)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

// TestCreateEntry tests CreateEntry func
func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)

	createRandomEntry(t, account)
}

// TestGetEntry tests GetEntry func
func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)

	entry1 := createRandomEntry(t, account)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.ID, entry2.ID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// TestListEntries tests ListEntries func
func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}

	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
