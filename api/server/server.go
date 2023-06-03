package server

import (
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func New(store *db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// receives the createAccount function
	server.router.POST("/accounts", server.createAccount)
	return server
}

// Start runs the server in the specified address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errResponse returns a gin key-value error
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
