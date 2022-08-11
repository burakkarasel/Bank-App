package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker which implements maker interface
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username for a specific duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// here we create a new payload for the token
	payload, err := NewPayload(username, duration)

	// check for error if any error occurs we return an empty string and the error
	if err != nil {
		return "", nil, err
	}

	// if no error occurs we create jwtToken with the payload we created for this token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(maker.secretKey))
	return signedToken, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// keyFunc will check given token's header if it matches or not
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// since we've used jwt.SigningMethodHS256 we check if header implements jwt.SigningMethodHMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// if it doesnt then given token doesnt match with our signingMethod
			return nil, ErrInvalidToken
		}
		// if it does then given token matches  with our signingMethod
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	// if an error occurs there are two possibilities first token might be invalid (might use another signing algorithm)
	// second it can be expired so we need to check the error for these both posibilities
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		// here we chek if any validation error occured during parsing the token
		// if it's exactly ErrExpiredToken we return nil Payload and ErrExpiredToken
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		// Otherwise we return nil Payload and ErrInvalidToken because our token is invalid
		return nil, ErrInvalidToken
	}

	// if everything is ok and no error occured we convert jwtToken Claims to Payload
	payload, ok := jwtToken.Claims.(*Payload)

	// and if an error occurs during converting we return invalidToken error
	if !ok {
		return nil, ErrInvalidToken
	}

	// and if everything goes fine (token is valid and not expired) we return the payload and no error
	return payload, nil
}
