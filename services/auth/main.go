package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

// Implement Login method
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Mock authentication logic
	if req.Username == "admin" && req.Password == "password" {
		return &pb.LoginResponse{
			Token:   "mock-token",
			Message: "Login successful",
		}, nil
	}
	return nil, status.Errorf(401, "Invalid username or password")
}

// Implement Register method
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Mock registration logic
	return &pb.RegisterResponse{
		UserId:  "mock-user-id",
		Message: "Registration successful",
	}, nil
}

func main() {
	port := os.Getenv("AUTH_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Println("Auth Service is running on port" + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
