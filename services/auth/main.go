package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/assidiqi598/umrah-erp/services/auth/db"
	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/assidiqi598/umrah-erp/services/auth/repositories"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	WaNumber string `json:"wa_number"`
	jwt.RegisteredClaims
}

// Implement Login method
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	repo := repositories.NewUserRepository()

	if req.Email != "" {
		// Fetch user from MongoDB
		user, err := repo.FindByEmail(req.Email)
		if err != nil {
			log.Printf("User not found: %v", err)
			return nil, status.Errorf(codes.NotFound, "User not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

		if err != nil {
			log.Printf("Invalid password for user: %v", user.Email)
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		}

		// Generate tokens
		accessToken, err := generateJWT(user.ID, user.Email, user.WhatsAppNumber, time.Minute*60)
		if err != nil {
			log.Printf("Error generating access token: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to generate access token")
		}

		refreshToken, err := generateJWT(user.ID, user.Email, user.WhatsAppNumber, time.Hour*24)
		if err != nil {
			log.Printf("Error generating refresh token: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to generate refresh token")
		}

		objectID, err := primitive.ObjectIDFromHex(user.ID)

		if err != nil {
			log.Printf("Error converting id: %v", err)
		}

		err = repo.UpdateUser(
			context.Background(),
			bson.M{"_id": objectID},
			bson.M{
				"$set": bson.M{
					"last_login": time.Now(),
				},
			})

		if err != nil {
			log.Printf("Error updating user: %v", err)
		}

		// Return a successful response
		return &pb.LoginResponse{Token: accessToken, Message: "Login successful", RefreshToken: refreshToken}, nil

	} else if req.WaNumber != "" {

		// Fetch user based on whatsapp number
		// user, err := repo.FindByWaNumber(req.WaNumber)
		// if err != nil {
		// 	log.Printf("User not found: %v", err)
		// 	return nil, status.Errorf(codes.NotFound, "User not found")
		// }

	}

	return nil, status.Error(codes.Unauthenticated, "Please provide the credentials")
}

func generateJWT(userID, email string, waNumber string, duration time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT secret key is not configured")
	}

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		WaNumber: waNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
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

	_, err = repo.FindByWaNumber(req.WaNumber)

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Whatsapp number already exists")
	}

	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking whatsapp number existence: %v", err)
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

	log.Println("Auth Service is running on port: " + port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
