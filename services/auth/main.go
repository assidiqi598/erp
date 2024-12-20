package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/assidiqi598/umrah-erp/services/auth/db"
	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/assidiqi598/umrah-erp/services/auth/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

// Implement Login method
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	repo := repositories.NewUserRepository()

	// Fetch user from MongoDB
	user, err := repo.FindByEmail(req.Email)
	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Check password (for simplicity, plaintext comparison is shown)
	if user.Password != req.Password {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	// Return a successful response
	return &pb.LoginResponse{Message: "Login successful"}, nil
}

// Implement Register method
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	repo := repositories.NewUserRepository()

	_, err := repo.FindByEmail(req.Email)

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Email already exists")
	}

	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking email existence: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		fmt.Println("Error:", err)
		cost = bcrypt.DefaultCost
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	// Insert the user
	newUser := repositories.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       string(hashedPassword),
		WhatsAppNumber: req.WaNumber,
		CreatedAt:      time.Now(),
	}

	err = repo.CreateUser(&newUser)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to create user")
	}

	return &pb.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	if os.Getenv("BCRYPT_COST") == "" {
		log.Fatal("BCRYPT_COST is not set")
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
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Println("Auth Service is running on port" + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
