package db

import (
	"context"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	args := CreateTransferParams{
		FromAccountID: createRandomAccount(t).ID,
		ToAccountID:   createRandomAccount(t).ID,
		Amount:        utils.RandomBalance(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, transfer)

	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	tr := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), tr.ID)
	require.NoError(t, err)

	require.NotEmpty(t, transfer)

	require.Equal(t, tr.Amount, transfer.Amount)
	require.Equal(t, tr.ToAccountID, transfer.ToAccountID)
	require.Equal(t, tr.FromAccountID, transfer.FromAccountID)

	require.WithinDuration(t, tr.CreatedAt.Time, transfer.CreatedAt.Time, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	tr := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), tr.ID)
	require.NoError(t, err)

	emptyTransfer, err := testQueries.GetTransfer(context.Background(), tr.ID)
	require.Error(t, err)
	require.Empty(t, emptyTransfer)
}

func TestGetTransferList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
