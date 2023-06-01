package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTxStore(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(1000)

	errs := make(chan error) // channels allow goroutines to communicate. connects concurrent goroutines.
	results := make(chan TransferTxResult)

	// implement go routines to avoid concurrency issues
	for i := 0; i < 5; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	for i := 0; i < 5; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// validate transfer
		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// validate get transfer
		transfer, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)

		// validate entries
		// fromEntry - transfer sender
		fromEntry := result.FromEntry
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		fromEntry, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NotEmpty(t, fromEntry)
		require.NoError(t, err)

		// toEntry - transfer receiver
		toEntry := result.ToEntry
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		fromEntry, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NotEmpty(t, toEntry)
		require.NoError(t, err)

	}
}
