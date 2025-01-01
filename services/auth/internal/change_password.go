package internal

import (
	"context"
	"log"
	"os"
	"strconv"

	pb "github.com/assidiqi598/erp/services/auth/proto"
	"github.com/assidiqi598/erp/shared/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {

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

	// Compare old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GivenPassword))

	if err != nil {
		log.Printf("Invalid password for user: %v", user.Email)
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		log.Printf("Error: %v", err)
		cost = bcrypt.DefaultCost
	}

	// Hash the new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), cost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	userObjectID, err := primitive.ObjectIDFromHex(user.ID)

	if err != nil {
		log.Printf("Error converting id: %v", err)
	}

	err = repo.UpdateUser(
		context.Background(),
		bson.M{"_id": userObjectID},
		bson.M{
			"$set": bson.M{
				"password": newHashedPassword,
			},
		})

	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	return &pb.ChangePasswordResponse{Message: "Password berhasil diubah"}, nil
}
