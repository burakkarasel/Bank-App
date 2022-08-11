package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/burakkarasel/Bank-App/api"
	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/gapi"
	"github.com/burakkarasel/Bank-App/pb"
	"github.com/burakkarasel/Bank-App/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".") // we pass the location of the file

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// here we connect to DB if any error occurs program shutsdown
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)
	// runGinServer(config, store)
}

// runGinServer runs a gin HTTP server
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// and here we start listening our server
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

// runGrpcServer runs a gRPC server
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBankAppServer(grpcServer, server)
	reflection.Register(grpcServer) // self doc

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("gRPC server started at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start server")
	}
}
