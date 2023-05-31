package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccountParams{
		Owner:    "Marina",
		Balance:  "100",
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)
}
