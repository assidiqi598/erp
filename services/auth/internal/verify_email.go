package internal

import (
	"context"
	"log"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	public "github.com/assidiqi598/umrah-erp/services/auth/public"
	"github.com/assidiqi598/umrah-erp/shared/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// Retrieve claims from the context
	claims, ok := ctx.Value(public.ClaimsKey).(*public.Claims)
	if !ok {
		log.Println("Failed to retrieve claims from context")
	}

	repo := repositories.NewUserRepository()

	userObjectID, err := primitive.ObjectIDFromHex(claims.UserID)

	if err != nil {
		log.Printf("Error converting id: %v", err)
	}

	// Fetch user from MongoDB
	user, err := repo.FindUser(bson.M{"_id": userObjectID})
	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if user.Token != req.EmailToken {
		log.Printf("Token for user %s not match", user.ID)
		return &pb.VerifyEmailResponse{Message: "Token salah"}, nil
	}

	err = repo.UpdateUser(
		context.Background(),
		bson.M{"_id": userObjectID},
		bson.M{
			"$set": bson.M{
				"is_verified": true,
			},
		})

	if err != nil {
		log.Printf("Error updating user verification of %s because %v", user.ID, err)
		return nil, status.Errorf(codes.Internal, "Error updating user")
	}

	return &pb.VerifyEmailResponse{Message: "Anda berhasil terverifikasi"}, nil
}
