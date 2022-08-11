package gapi

import (
	"errors"
	"fmt"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/pb"
	"github.com/burakkarasel/Bank-App/token"
	"github.com/burakkarasel/Bank-App/util"
)

var ErrAccountIsNotAuthenticatedUsers = errors.New("account doesn't belong to authenticated user")

// Server serves all HTTP request for banking services
type Server struct {
	pb.UnimplementedBankAppServer
	config     util.Config
	store      db.Store // which we will hold the db, and queries
	tokenMaker token.Maker
}

// NewServer creates a new Server which will hold our routes and DB
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
