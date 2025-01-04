package main

import (
	"log"
	"net"
	"os"

	"github.com/assidiqi598/erp/gateway/internal"
	authpb "github.com/assidiqi598/erp/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	ensureENV()

	// Load TLS credentials for Flutter ⇔ Gateway
	gatewayCreds, err := credentials.NewServerTLSFromFile(os.Getenv("CERT_FILE_GW"), os.Getenv("KEY_FILE"))
	if err != nil {
		log.Fatalf("Failed to load gateway TLS credentials: %v", err)
	}

	// Load TLS credentials for Gateway ⇔ Microservices
	microservicesCreds, err := credentials.NewClientTLSFromFile(os.Getenv("CERT_FILE_MS"), "")
	if err != nil {
		log.Fatalf("Failed to load microservices TLS credentials: %v", err)
	}

	// Create clients for microservices
	authConn, err := grpc.NewClient(os.Getenv("AUTH_SERVICE"), grpc.WithTransportCredentials(microservicesCreds))
	if err != nil {
		log.Fatalf("Failed to connect to Auth service: %v", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	// Start the gRPC Gateway Server
	server := grpc.NewServer(grpc.Creds(gatewayCreds), grpc.UnaryInterceptor(internal.ForwardMetadataInterceptor()))
	authpb.RegisterAuthServiceServer(server, internal.NewAuthHandler(authClient))

	listener, err := net.Listen("tcp", "0.0.0.0:"+os.Getenv("GATEWAY_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen on 0.0.0.0:5000: %v", err)
	}

	reflection.Register(server)

	log.Printf("Gateway server is running on port %s..", os.Getenv("GATEWAY_PORT"))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func ensureENV() {
	if os.Getenv("GATEWAY_PORT") == "" {
		log.Fatal("GATEWAY_PORT is not set")
	}

	if os.Getenv("AUTH_SERVICE") == "" {
		log.Fatal("AUTH_SERVICE is not set")
	}

	if os.Getenv("CERT_FILE_MS") == "" {
		log.Fatal("CERT_FILE_MS is not set")
	}

	if os.Getenv("CERT_FILE_GW") == "" {
		log.Fatal("CERT_FILE_GW is not set")
	}

	if os.Getenv("KEY_FILE") == "" {
		log.Fatal("KEY_FILE is not set")
	}
}
