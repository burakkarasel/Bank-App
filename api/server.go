package api

import (
	"fmt"

	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/token"
	"github.com/burakkarasel/Bank-App/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all HTTP request for banking services
type Server struct {
	config     util.Config
	store      db.Store // which we will hold the db, and queries
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new Server which will hold our routes and DB
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// setupRouter holds our routes
func (server *Server) setupRouter() {
	router := gin.Default()

	// users no middleware for these routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// accounts
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountById)
	authRoutes.GET("/accounts", server.listAccounts)

	// transfers
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific port to handler requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse lets us to send error to client in JSON format (key:value)
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
