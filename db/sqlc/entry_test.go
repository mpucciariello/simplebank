package db

import (
	"context"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T) Entry {
	args := CreateEntryParams{
		AccountID: CreateRandomAccount(t).ID,
		Amount:    utils.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	e := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), e.ID)
	require.NoError(t, err)

	require.NotEmpty(t, entry)

	require.Equal(t, e.Amount, entry.Amount)
	require.Equal(t, e.AccountID, entry.AccountID)
	require.Equal(t, e.ID, entry.ID)

	require.WithinDuration(t, e.CreatedAt.Time, entry.CreatedAt.Time, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	e := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), e.ID)
	require.NoError(t, err)

	emptyEntry, err := testQueries.GetEntry(context.Background(), e.ID)
	require.Error(t, err)
	require.Empty(t, emptyEntry)
}

func TestGetEntryList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, account := range entries {
		require.NotEmpty(t, account)
	}
}
