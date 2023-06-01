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
	amount := int64(10)

	n := 5
	exists := make(map[int]bool)

	errs := make(chan error) // channels allow goroutines to communicate. connects concurrent goroutines.
	results := make(chan TransferTxResult)

	// implement go routines to avoid concurrency issues
	for i := 0; i < n; i++ {
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
	// n defines the number of executions
	for i := 0; i < n; i++ {

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

		// check accounts - from account
		fromAccount := result.FromAccountID
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.Owner, fromAccount.Owner)

		// check accounts - to account
		toAccount := result.ToAccountID
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		require.Equal(t, account2.Owner, toAccount.Owner)

		// check accounts balance
		diffAccount := account2.Balance - account1.Balance
		diffResult := toAccount.Balance - fromAccount.Balance

		require.Equal(t, diffAccount, diffResult)
		require.True(t, diffResult > 0)
		require.True(t, diffResult%amount == 0)

		// n is the number of times the transfer has been done
		z := int(diffResult / amount)
		require.True(t, z > 0 && z <= n)

		// k must be unique if concurrency is working properly
		require.NotContains(t, exists, z)
		exists[z] = true

		// check updated accounts
		updatedFromAccount, err := store.GetAccount(context.Background(), fromAccount.ID)
		require.NotEmpty(t, updatedFromAccount)
		require.NoError(t, err)

		updatedToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
		require.NotEmpty(t, updatedToAccount)
		require.NoError(t, err)

		require.Equal(t, account1.Balance-int64(z)*amount, updatedFromAccount.Balance)
		require.Equal(t, account2.Balance+int64(z)*amount, updatedToAccount.Balance)
	}
}
