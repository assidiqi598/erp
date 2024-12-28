package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestMain(m *testing.M) {
	// Load environment variables from .env file
	err := godotenv.Load(".env") // Adjust the path to your .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestAuthServiceE2EWithDB(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Verify the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to verify connection: %v", err)
	}

	defer mongoClient.Disconnect(ctx)

	usersCollection := mongoClient.Database(os.Getenv("DB_NAME")).Collection("users")

	// gRPC connection setup
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	const email = "developmore@yahoo.com"
	var userId string
	var emailToken string

	// Test Register
	t.Run("Register", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Email:       email,
			Password:    "password",
			Username:    "Test User",
			PhoneNumber: "085925119040",
		}

		res, err := client.Register(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Anda berhasil terdaftar, mohon verifikasi dengan token yang telah dikirim.")
		assert.NotNil(t, res.UserId)
		userId = res.UserId

		// Check database for new user
		var user bson.M
		err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, email, user["email"])
		assert.Equal(t, "Test User", user["username"])
		assert.Contains(t, user["verification_msg_id"], "mailin.fr")
		emailToken = user["token"].(string)
	})

	// Test Login
	t.Run("Login", func(t *testing.T) {
		req := &pb.LoginWithEmailAndPassRequest{
			Email:    email,
			Password: "password",
		}

		res, err := client.LoginWithEmailAndPass(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Login successful")
		assert.NotEmpty(t, res.Token, "Expected a non-empty JWT token")
		assert.NotEmpty(t, res.RefreshToken, "Expected a non-empty JWT refresh token")
	})

	// Test VerifyEmail (optional, depending on your service implementation)
	t.Run("VerifyEmail", func(t *testing.T) {
		req := &pb.VerifyEmailRequest{
			Email: email,
			Token: emailToken, // Replace with the expected verification code
		}

		res, err := client.VerifyEmail(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Expected email verification to be successful")

		// Optionally check the database for updated verification status
		var user bson.M
		err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, true, user["is_verified"])
	})

	t.Cleanup(func() {
		objectID, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			fmt.Printf("Error converting id: %v", err)
		}
		usersCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	})
}
