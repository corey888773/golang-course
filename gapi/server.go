package gapi

import (
	"fmt"

	db "github.com/corey888773/golang-course/db/sqlc"
	"github.com/corey888773/golang-course/pb"
	"github.com/corey888773/golang-course/token"
	"github.com/corey888773/golang-course/util"
	"github.com/corey888773/golang-course/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

// new gRPC server
func NewServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a maker: %w", err)
	}

	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
