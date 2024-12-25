package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/assidiqi598/umrah-erp/services/auth/db"
	"github.com/assidiqi598/umrah-erp/services/auth/internal"
	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	if os.Getenv("BCRYPT_COST") == "" {
		log.Fatal("BCRYPT_COST is not set")
	}

	if os.Getenv("BREVO_API_KEY") == "" {
		log.Fatal("BREVO_API_KEY is not set")
	}

	// Connect to MongoDB
	err := db.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Client.Disconnect(context.TODO())

	db.CreateUniqueIndex()

	port := os.Getenv("AUTH_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &internal.AuthServer{})

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Println("Auth Service is running on port: " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
