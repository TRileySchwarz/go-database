package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/TRileySchwarz/go-database/models"
	"github.com/dgrijalva/jwt-go"
)

// The secret JWT token key used
var jwtKey = []byte("tokenKey")

// This represents the amount of time a JWT is valid for
var authDuration = time.Duration(10 * time.Minute)

// Generates a new JWT using the userEmail / ID
func GenerateJWT(userEmail string) (string, error) {
	timeLimit := time.Now().Add(authDuration)
	claims := &models.Claims{
		Username: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: timeLimit.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Checks that a json webToken is valid, and returns the corresponding email if valid
func VerifyWebToken(webToken string) (string, error) {
	claims := &models.Claims{}

	// Verify the claim is valid and extract the email/ID
	token, err := jwt.ParseWithClaims(webToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return " ", err
		}
		return " ", err
	}
	// If token is invalid throw an error
	if !token.Valid {
		return " ", errors.New("invalid token has been sent")
	}

	return claims.Username, nil
}

// HashPassword converts a user password into hashed version stored in the db
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

