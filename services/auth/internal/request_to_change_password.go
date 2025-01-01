package internal

import (
	"context"
	"log"
	"os"
	"strconv"

	pb "github.com/assidiqi598/erp/services/auth/proto"
	"github.com/assidiqi598/erp/shared/repositories"
	"github.com/assidiqi598/erp/shared/storage"
	"github.com/assidiqi598/erp/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) RequestToChangePassword(
	ctx context.Context,
	req *pb.RequestToChangePasswordRequest,
) (*pb.RequestToChangePasswordResponse, error) {

	repo := repositories.NewUserRepository()

	// Fetch user from MongoDB
	user, err := repo.FindUser(bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"phone_number": req.PhoneNumber},
		},
	})
	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	randomString, err := utils.GenerateSecureRandomString(10)

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		log.Printf("Error: %v", err)
		cost = bcrypt.DefaultCost
	}

	// Hash the random string
	hashedRandomString, err := bcrypt.GenerateFromPassword([]byte(randomString), cost)

	userObjectId, err := primitive.ObjectIDFromHex(user.ID)

	if err != nil {
		log.Printf("Error converting id: %v", err)
	}

	err = repo.UpdateUser(
		context.Background(),
		bson.M{"_id": userObjectId},
		bson.M{
			"$set": bson.M{
				"given_password": hashedRandomString,
			},
		})

	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	emailData := struct {
		Username      string
		GivenPassword string
	}{
		Username:      user.Username,
		GivenPassword: randomString,
	}

	s3Client := storage.GetS3Client()

	emailHTML, err := s3Client.GetEmailTemplateAndReplace(
		os.Getenv("S3_BUCKET_NAME"),
		"email_templates/change_password.html",
		emailData,
	)

	if err != nil {
		log.Printf("Error getting html email content: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed getting html email content")
	}

	resp, err := utils.SendEmail(
		os.Getenv("BREVO_API_KEY"),
		"do-not-reply@devmore.id",
		"Devmore",
		user.Email,
		user.Username,
		"Verifikasi Email",
		"Berikut merupakan kode verifikasi email Anda",
		emailHTML,
	)

	if err != nil {
		log.Printf("Error sending email verification: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to send email verification")
	}

}
