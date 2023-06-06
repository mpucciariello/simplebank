package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func New(store db.Store) (server *Server) {
	server = &Server{store: store}

	router := gin.Default()
	server.initRouter(router)

	server.router = router
	return
}

// Start runs the server in the specified address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) initRouter(router *gin.Engine) {
	// declares the api routes and its functions
	router.POST("/accounts", s.createAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.GET("/accounts", s.getAccountsList)
	router.DELETE("/accounts/:id", s.deleteAccount)
}

// errResponse returns a gin key-value error
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
