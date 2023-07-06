package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/token"
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
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token validator: %w", err)
	}

	server = &Server{
		store:  store,
		token:  tokenMaker,
		config: config,
	}

	// set currency validator
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
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	authRoutes := router.Group("/", authMiddleware(s.token))
	authRoutes.GET("/users/:username", s.getUser)

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.getAccountsList)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTranfer)
}

// errResponse returns a gin key-value error
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
