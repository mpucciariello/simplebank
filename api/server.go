package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (server *Server) {
	server = &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}

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

	router.POST("/transfers", s.createTranfer)
}

// errResponse returns a gin key-value error
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}