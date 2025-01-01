package internal

import (
	"context"
	"log"
	"os"

	pb "github.com/assidiqi598/erp/services/auth/proto"
	"github.com/assidiqi598/erp/services/auth/public"
	"github.com/assidiqi598/erp/shared/repositories"
	"github.com/assidiqi598/erp/shared/storage"
	"github.com/assidiqi598/erp/shared/utils"
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
		return nil, status.Errorf(codes.Internal, "Gagal mengidentifikasi request.")
	}

	repo := repositories.NewUserRepository()

	userObjectId, err := primitive.ObjectIDFromHex(claims.UserID)

	if err != nil {
		log.Printf("Error converting id: %v", err)
		return nil, status.Errorf(codes.Internal, "Terjadi kesalahan konversi ID.")
	}

	user, err := repo.FindUser(bson.M{"_id": userObjectId})

	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "User tidak ditemukan.")
	}

	s3Client := storage.GetS3Client()

	emailHTML, err := s3Client.GetEmailTemplateAndReplace(
		os.Getenv("S3_BUCKET_NAME"),
		"email_templates/verifikasi_token.html",
		user,
	)

	if err != nil {
		log.Printf("Error getting html email content: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal menyiapkan email.")
	}

	resp, err := utils.SendEmail(
		os.Getenv("BREVO_API_KEY"),
		"do-not-reply@devmore.id",
		"Devmore",
		user.Email,
		user.Username,
		"Verifikasi Email",
		"Berikut merupakan kode verifikasi email Anda.",
		emailHTML,
	)

	if err != nil {
		log.Printf("Error sending email verification: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal mengirim email verifikasi.")
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
		return nil, status.Errorf(codes.Internal, "Gagal update user.")
	}

	return &pb.ResendVerificationEmailResponse{Message: "Email telah terkirim kembali."}, nil
}
