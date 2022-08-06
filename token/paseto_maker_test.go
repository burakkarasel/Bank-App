package token

import (
	"testing"
	"time"

	"github.com/burakkarasel/Bank-App/util"
	"github.com/stretchr/testify/require"
)

// TestPasetoMaker tests PasetoMaker func
func TestPasetoMaker(t *testing.T) {
	// we create a new maker with random 32 characters string
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	// then we create a new random owner for username
	username := util.RandomOwner()
	duration := time.Minute

	// we assign time.Now to issued at, and add duration to it and assign it to expired at
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	// then we create token with these information
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// and then we verify the token and create a payload with given informations
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// TestExpiredPasetoToken tests an expired Paseto Token
func TestExpiredPasetoToken(t *testing.T) {
	// we create a new maker
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	// then we create a new token with -1 minute duration which will be always expired
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// then we verify the token and it should return expired error
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
