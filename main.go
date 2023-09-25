package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/corey888773/golang-course/api"
	db "github.com/corey888773/golang-course/db/sqlc"
	"github.com/corey888773/golang-course/gapi"
	"github.com/corey888773/golang-course/pb"
	"github.com/corey888773/golang-course/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:")
	}

	connector, err := pq.NewConnector(config.DbSource)
	if err != nil {
		log.Fatal("cannot create connector:", err)
	}
	connection := sql.OpenDB(connector)

	defer connection.Close()

	store := db.NewStore(connection)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	grpcServer := grpc.NewServer()
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
