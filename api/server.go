package api

import (
	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all HTTP request for banking services
type Server struct {
	store  db.Store // which we will hold the db, and queries
	router *gin.Engine
}

// NewServer creates a new Server which will hold our routes and DB
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// accounts
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)

	// transfers
	router.POST("/transfers", server.createTransfer)

	// users
	router.POST("/users", server.createUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific port to handler requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse lets us to send error to client in JSON format (key:value)
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func successResponse(msg string) gin.H {
	return gin.H{"success": msg}
}
