package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestPassword tests our password funcs
func TestPassword(t *testing.T) {
	// test ok case
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// test fail case
	wrongPassword := RandomString(8)

	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// check if two other hash for same password is different
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
