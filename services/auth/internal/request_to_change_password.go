package internal

import (
	"context"
	"log"

	pb "github.com/assidiqi598/umrah-erp/services/auth/proto"
	"github.com/assidiqi598/umrah-erp/shared/repositories"
	"go.mongodb.org/mongo-driver/bson"
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

}
