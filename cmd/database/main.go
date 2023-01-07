package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.fsrv.services/fsrvcorp/analytics/tracking-server/pkg/database"
	"gitlab.fsrv.services/fsrvcorp/analytics/tracking-server/pkg/proto"
)

type server struct {
	proto.UnimplementedIngestServiceServer
	database *database.Settings
}

func (s server) IngestMessage(ctx context.Context, request *proto.Request) (*proto.IngestResponse, error) {
	mad := make(map[string]interface{})
	for _, v := range request.GetHeaders() {
		mad[v.Key] = v.Value
	}

	tx := s.database.Client.WithContext(ctx).Create(&database.Request{
		Headers: mad,
		Path:    request.GetPath(),
	})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &proto.IngestResponse{}, nil
}

func main() {
	db := database.Settings{
		Host:     os.Getenv("DATABASE_HOST"),
		Username: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}
	databaseInitializationError := db.InitializeDB(log.New(os.Stdout, "", log.LstdFlags))
	if databaseInitializationError != nil {
		log.Fatalf("can not initialize database %v", databaseInitializationError)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	proto.RegisterIngestServiceServer(s, &server{
		database: &db,
	})

	log.Println("start server")
	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
