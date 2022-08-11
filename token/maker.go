package token

import "time"

// Maker is an interface for managing tokens so we can change between JWT & PASETO
type Maker interface {
	// CreateToken creates a new token for a specific username for a specific duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not, if its valid VerifyToken method will return the payload of token
	VerifyToken(token string) (*Payload, error)
}
