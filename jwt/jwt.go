package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtSecret = "aJWTSecret"

// CreateJWTToken expects user ID and username and creates a JWT token
func CreateJWTToken(userID, userName string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // expiry: one week
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

// GetUserFromJWTToken expects JWT token, decodes the user ID and username
func GetUserFromJWTToken(token *jwt.Token) (string, string, error) {
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	userName := claims["username"].(string)
	if userID == "" || userName == "" {
		return "", "", fmt.Errorf("user ID or user name is empty")
	}

	return userID, userName, nil
}
