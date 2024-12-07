package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Payload struct for the JWT token
type Payload struct {
	UserID string
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a given userID
func GenerateJWT(userID string, tokenExpiry time.Duration) (string, error) {
	// Create JWT Payload
	payload := Payload{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiry).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte("mysecret"))
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing token")
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return "", fmt.Errorf("invalid token format")
}
