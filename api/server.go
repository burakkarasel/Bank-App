package api

import (
	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves all HTTP request for banking services
type Server struct {
	store  *db.Store // which we will hold the db, and queries
	router *gin.Engine
}

// NewServer creates a new Server which will hold our routes and DB
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.listAccounts)

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
