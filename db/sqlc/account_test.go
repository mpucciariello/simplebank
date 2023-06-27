package db

import (
	"context"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.UserName,
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	a := CreateRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), a.ID)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, a.Owner, account.Owner)
	require.Equal(t, a.Balance, account.Balance)
	require.Equal(t, a.Currency, account.Currency)
	require.Equal(t, a.ID, account.ID)

	require.WithinDuration(t, a.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	a := CreateRandomAccount(t)

	args := UpdateAccountParams{
		ID:      a.ID,
		Balance: utils.RandomBalance(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	require.Equal(t, a.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, a.Currency, account.Currency)
	require.Equal(t, a.ID, account.ID)

	require.WithinDuration(t, a.CreatedAt.Time, account.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	a := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), a.ID)
	require.NoError(t, err)

	emptyAccount, err := testQueries.GetAccount(context.Background(), a.ID)
	require.Error(t, err)
	require.Empty(t, emptyAccount)
}

func TestGetAccountList(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
