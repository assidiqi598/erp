package test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/assidiqi598/erp/services/auth/proto"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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

	creds, err := credentials.NewClientTLSFromFile(os.Getenv("CERT_FILE"), "")
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}

	// gRPC connection setup
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	const email = "developmore@yahoo.com"
	var emailToken string
	var loginToken string
	var verificationMsgId string

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
		assert.Equal(t, res.Message, "Anda berhasil terdaftar, mohon login dan verifikasi dengan token yang telah dikirim.")
		assert.NotNil(t, res.UserId)

		// Check database for new user
		var user bson.M
		err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, email, user["email"])
		assert.Equal(t, "Test User", user["username"])
		assert.Contains(t, user["verification_msg_id"], "mailin.fr")

		emailToken = user["email_token"].(string)
		verificationMsgId = user["verification_msg_id"].(string)
	})

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

		loginToken = res.Token
	})

	t.Run("ResendVerificationEmail", func(t *testing.T) {
		md := metadata.New(map[string]string{"authorization": "Bearer " + loginToken})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		req := &pb.ResendVerificationEmailRequest{}

		res, err := client.ResendVerificationEmail(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Email telah terkirim kembali.")

		var user bson.M
		err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		assert.NoError(t, err)
		assert.Contains(t, user["verification_msg_id"], "mailin.fr")
		assert.NotSame(t, verificationMsgId, user["verification_msg_id"])
	})

	t.Run("VerifyEmail", func(t *testing.T) {

		md := metadata.New(map[string]string{"authorization": "Bearer " + loginToken})
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		req := &pb.VerifyEmailRequest{
			EmailToken: emailToken, // Replace with the expected verification code
		}

		res, err := client.VerifyEmail(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Message, "Anda berhasil terverifikasi")

		// Optionally check the database for updated verification status
		var user bson.M
		err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		assert.NoError(t, err)
		assert.Equal(t, true, user["is_verified"])
	})

	// t.Run("Remove Test User", func(t *testing.T) {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()

	// 	res, err := usersCollection.DeleteOne(ctx, bson.M{"email": email})
	// 	if err != nil {
	// 		t.Fatalf("Failed to clean up test user: %v", err)
	// 	}

	// 	log.Printf("Cleaned up test user with email: %s", email)

	// 	assert.Equal(t, 1, int(res.DeletedCount))

	// })
}
