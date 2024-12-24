package repositories

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/assidiqi598/umrah-erp/services/auth/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID          string    `bson:"_id,omitempty"`
	Username    string    `bson:"username"`
	Password    string    `bson:"password"`
	Email       string    `bson:"email"`
	PhoneNumber string    `bson:"phone_number"`
	IsVerified  bool      `bson:"is_verified"`
	Token       string    `bson:"token"`
	LastLogin   time.Time `bson:"last_login"`
	CreatedAt   time.Time `bson:"created_at"`
}

// UserRepository provides methods for user-related database operations
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: db.GetCollection(os.Getenv("DB_NAME"), "users"),
	}
}

// FindUser fetches a user by a filter
func (r *UserRepository) FindUser(filter bson.M) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// UpdateUser updates a user based on filter and update
func (r *UserRepository) UpdateUser(ctx context.Context, filter bson.M, update bson.M) error {
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
