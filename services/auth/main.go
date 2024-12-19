package main

import (
	"context"
	"log"
	"net"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"google.golang.org/grpc"
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
	return nil, grpc.Errorf(401, "Invalid username or password")
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
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	log.Println("Auth Service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
