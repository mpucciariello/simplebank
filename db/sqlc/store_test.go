package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTxStore(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> before transfer:", account1.Balance, account2.Balance)
	amount := int64(10)

	n := 10
	exists := make(map[int]bool)

	errs := make(chan error) // channels allow goroutines to communicate. connects concurrent goroutines.
	results := make(chan TransferTxResult)

	// implement go routines to avoid concurrency issues
	for i := 0; i < n; i++ {
		trxName := fmt.Sprintf("trx: %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, trxName)

			result, err := store.TransferTx(ctx, TransferTxParams{
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
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff2 > 0)
		require.True(t, diff2%amount == 0)

		// n is the number of times the transfer has been done
		z := int(diff2 / amount)
		require.True(t, z > 0 && z <= n)

		// k must be unique if concurrency is working properly
		require.NotContains(t, exists, z)
		exists[z] = true

	}

	// check updated accounts
	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedFromAccount)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedToAccount)
	require.NoError(t, err)

	fmt.Println(">> after transfer:", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTxStoreDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> before transfer:", account1.Balance, account2.Balance)
	amount := int64(10)

	n := 10 // 5 from account1 to account 2, 5 from account2 to account1

	errs := make(chan error)

	for i := 0; i < n; i++ {
		trxName := fmt.Sprintf("trx: %d", i+1)

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, trxName)

			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	//check results
	// n defines the number of executions
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check updated accounts
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)

	fmt.Println(">> after transfer:", updatedAccount1.Balance, updatedAccount2.Balance)

	// should be equal since at the end the balance is the same as before testing
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
