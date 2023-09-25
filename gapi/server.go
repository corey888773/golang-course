package gapi

import (
	"fmt"

	db "github.com/corey888773/golang-course/db/sqlc"
	"github.com/corey888773/golang-course/pb"
	"github.com/corey888773/golang-course/token"
	"github.com/corey888773/golang-course/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

// new gRPC server
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
