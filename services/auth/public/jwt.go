package internal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email string, phoneNumber string, duration time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT secret key is not configured")
	}

	claims := &Claims{
		UserID:      userID,
		Email:       email,
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// Function to extract userID from JWT
func extractUserIDFromToken(secretKey string, tokenString string) (string, error) {
	// Parse the JWT string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the userID from claims
		if userID, ok := claims["UserID"].(string); ok {
			return userID, nil
		}
		return "", fmt.Errorf("userID not found in the token")
	}
	return "", fmt.Errorf("invalid token")
}
