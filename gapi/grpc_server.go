package gapi

import (
	"fmt"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/pb"
	"github.com/micaelapucciariello/simplebank/token"
	"github.com/micaelapucciariello/simplebank/utils"
)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedSimpleBankServer
	store  db.Store
	token  token.Maker
	config utils.Config
}

func NewServer(config utils.Config, store db.Store) (server *Server, err error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token validator: %w", err)
	}

	server = &Server{
		store:  store,
		token:  tokenMaker,
		config: config,
	}

	return
}
