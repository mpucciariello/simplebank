package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/micaelapucciariello/simplebank/api/token"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	token  token.Maker
	config utils.Config
}

func NewServer(config utils.Config, store db.Store) (server *Server, err error) {
	router := gin.Default()
	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token validator: %w", err)
	}

	server = &Server{store: store, token: token, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
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

	router.POST("/users", s.createUser)
	router.GET("/users/:username", s.getUser)
	router.POST("/users/login", s.loginUser)
}

// errResponse returns a gin key-value error
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
