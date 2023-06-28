package db

import (
	"context"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	args := CreateUserParams{
		UserName:       utils.RandomOwner(),
		HashedPassword: "password",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, args.UserName, user.UserName)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	u := CreateRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), u.UserName)
	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, u.UserName, user.UserName)
	require.Equal(t, u.FullName, user.FullName)
	require.Equal(t, u.Email, user.Email)

	require.WithinDuration(t, u.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	u := CreateRandomUser(t)

	args := UpdateUserParams{
		UserName: u.UserName,
		Email:    utils.RandomEmail(),
		FullName: utils.RandomOwner(),
	}

	user, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)

	require.Equal(t, args.UserName, user.UserName)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotEmpty(t, user)

	require.WithinDuration(t, u.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	u := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), u.UserName)
	require.NoError(t, err)

	emptyAccount, err := testQueries.GetUser(context.Background(), u.UserName)
	require.Error(t, err)
	require.Empty(t, emptyAccount)
}
