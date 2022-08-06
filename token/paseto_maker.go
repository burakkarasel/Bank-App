package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username for a specific duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// we create a new payload
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	// then we call this func and pass it our symmetricKey, the payload we created, and no footer
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// here we decrypt the token to verify into payload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}

	// after checking for invalid token we check if token is expired or not here
	err = payload.Valid()

	// valid returns expiredToken error no need to specify again
	if err != nil {
		return nil, err
	}

	return payload, nil
}
