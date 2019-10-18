package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"

	"github.com/TRileySchwarz/go-database/models"
	"github.com/dgrijalva/jwt-go"
)

// The secret JWT token key used
var jwtKey = []byte("tokenKey")

// This represents the amount of time a JWT is valid for
var authDuration = 10 * time.Minute

// Indicates the minimum password strength for a pass to be considered valid
var MinPassStrength = 5

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


// The password complexity functions are derived from the following library:
// https://github.com/briandowns/GoPasswordUtilities
type Password struct {
	Pass            string
	Length          int
	Score           int
	ContainsUpper   bool
	ContainsLower   bool
	ContainsNumber  bool
	ContainsSpecial bool
	ContainsLength bool
}

// ProcessPassword will parse the password and populate the Password struct attributes
// this helps check the strength of a potential password
func ProcessPassword(password string) *Password {
	p := &Password{Pass: password, Length: len(password)}

	matchLower := regexp.MustCompile(`[a-z]`)
	matchUpper := regexp.MustCompile(`[A-Z]`)
	matchNumber := regexp.MustCompile(`[0-9]`)
	matchSpecial := regexp.MustCompile(`[\!\@\#\$\%\^\&\*\(\\\)\-_\=\+\,\.\?\/\:\;\{\}\[\]~]`)

	if p.Length >= 8 {
		p.ContainsLength = true
		p.Score++
	}
	if matchLower.MatchString(p.Pass) {
		p.ContainsLower = true
		p.Score++
	}
	if matchUpper.MatchString(p.Pass) {
		p.ContainsUpper = true
		p.Score++
	}
	if matchNumber.MatchString(p.Pass) {
		p.ContainsNumber = true
		p.Score++
	}
	if matchSpecial.MatchString(p.Pass) {
		p.ContainsSpecial = true
		p.Score++
	}

	return p
}



