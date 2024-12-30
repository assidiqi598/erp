package internal

import (
	"context"
	"log"
	"os"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/assidiqi598/umrah-erp/services/auth/public"
	"github.com/assidiqi598/umrah-erp/shared/repositories"
	"github.com/assidiqi598/umrah-erp/shared/storage"
	"github.com/assidiqi598/umrah-erp/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) ResendVerificationEmail(
	ctx context.Context,
	req *pb.ResendVerificationEmailRequest,
) (*pb.ResendVerificationEmailResponse, error) {
	// Retrieve claims from the context
	claims, ok := ctx.Value(public.ClaimsKey).(*public.Claims)
	if !ok {
		log.Println("Failed to retrieve claims from context")
	}

	repo := repositories.NewUserRepository()

	userObjectId, err := primitive.ObjectIDFromHex(claims.UserID)

	if err != nil {
		log.Printf("Error converting id: %v", err)
	}

	user, err := repo.FindUser(bson.M{"_id": userObjectId})

	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	emailHTML, err := storage.GetEmailTemplateAndReplace(os.Getenv("S3_URI")+"email_templates/verifikasi_token.html", user)

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

	err = repo.UpdateUser(
		context.Background(),
		bson.M{"_id": userObjectId},
		bson.M{
			"$set": bson.M{
				"verification_msg_id": resp,
			},
		})

	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	return &pb.ResendVerificationEmailResponse{Message: "Email telah terkirim kembali"}, nil
}
