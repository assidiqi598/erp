package internal

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/assidiqi598/umrah-erp/shared/repositories"
	"github.com/assidiqi598/umrah-erp/shared/storage"
	"github.com/assidiqi598/umrah-erp/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implement Register method
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	repo := repositories.NewUserRepository()

	// Find user by email
	_, err := repo.FindUser(bson.M{"email": req.Email})

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Email already exists")
	}

	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking email existence: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	// Find user by phone number
	_, err = repo.FindUser(bson.M{"phone_number": req.PhoneNumber})

	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Phone number already exists")
	}

	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking phone number existence: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		log.Printf("Error: %v", err)
		cost = bcrypt.DefaultCost
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	token := strconv.Itoa(utils.GenerateSecureRandomNumber(6))

	// Insert the user
	newUser := repositories.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Token:       token,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
	}

	emailHTML, err := storage.GetEmailTemplateAndReplace(os.Getenv("S3_URI")+"email_templates/verifikasi_token.html", newUser)

	if err != nil {
		log.Printf("Error getting html email content: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed getting html email content")
	}

	resp, err := utils.SendEmail(os.Getenv("BREVO_API_KEY"), "do-not-reply@devmore.id", "Devmore", req.Email, req.Username, "Verifikasi Email", "Berikut merupakan kode verifikasi email Anda", emailHTML)

	if err != nil {
		log.Printf("Error sending email verification: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to send email verification")
	}

	newUser.VerificationMsgID = resp

	err = repo.CreateUser(&newUser)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to create user")
	}

	return &pb.RegisterResponse{
		Message: "Anda berhasil terdaftar, mohon verifikasi dengan token yang telah dikirim.",
	}, nil
}
