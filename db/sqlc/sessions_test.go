package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomSession(t *testing.T) Session {
	args := CreateSessionParams{
		ID:           uuid.New(),
		Username:     utils.RandomString(6),
		RefreshToken: utils.RandomString(32),
		UserAgent:    utils.RandomString(6),
		ClientIp:     utils.RandomString(6),
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{},
	}

	session, err := testQueries.CreateSession(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, session)

	require.Equal(t, args.Username, session.Username)

	require.NotZero(t, session.CreatedAt)
	require.NotZero(t, session.ID)

	return session
}

func TestCreateRandomSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	s := createRandomSession(t)

	session, err := testQueries.GetSession(context.Background(), s.ID)
	require.NoError(t, err)

	require.NotEmpty(t, session)

	require.Equal(t, s.Username, session.Username)
	require.Equal(t, s.RefreshToken, session.RefreshToken)
	require.Equal(t, s.UserAgent, session.UserAgent)
	require.Equal(t, s.ClientIP, session.ClientIP)
	require.Equal(t, s.ID, s.ID)

	require.WithinDuration(t, s.CreatedAt.Time, session.CreatedAt.Time, time.Second)
	require.WithinDuration(t, s.ExpiresAt.Time, session.ExpiresAt.Time, time.Second)
}
