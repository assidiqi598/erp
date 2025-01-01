package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/assidiqi598/erp/services/auth/internal"
	pb "github.com/assidiqi598/erp/services/auth/proto"
	"github.com/assidiqi598/erp/services/auth/public"
	"github.com/assidiqi598/erp/shared/db"
	"github.com/assidiqi598/erp/shared/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	if os.Getenv("S3_ENDPOINT") == "" {
		log.Fatal("S3_ENDPOINT is not set")
	}

	if os.Getenv("S3_ACCESS_KEY") == "" {
		log.Fatal("S3_ACCESS_KEY is not set")
	}

	if os.Getenv("S3_SECRET_KEY") == "" {
		log.Fatal("S3_SECRET_KEY is not set")
	}

	if os.Getenv("S3_BUCKET_NAME") == "" {
		log.Fatal("S3_BUCKET_NAME is not set")
	}

	// Connect to MongoDB
	err := db.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Client.Disconnect(context.TODO())

	db.CreateUniqueIndex()

	err = storage.CreateS3Client(storage.S3Credentials{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
	})

	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	port := os.Getenv("AUTH_PORT")

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Load the certificates into credentials
	creds, err := credentials.NewServerTLSFromFile(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
	if err != nil {
		log.Fatalf("failed to load TLS certificates: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(public.JwtAuthInterceptor))
	pb.RegisterAuthServiceServer(grpcServer, &internal.AuthServer{})

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Println("Auth Service is running on port: " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
