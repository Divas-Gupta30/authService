package helper

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims struct for the JWT token
type Claims struct {
	UserID string
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a given userID
func GenerateJWT(userID string) (string, error) {
	// Create JWT claims
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("mysecret"))
}
