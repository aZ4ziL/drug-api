package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secretkey")

func EncryptionPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func DecryptionPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

type Credential struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	LastLogin  time.Time `json:"last_login"`
	DateJoined time.Time `json:"date_joined"`
	IsAdmin    bool      `json:"is_admin"`
}

type Claims struct {
	Credential
	jwt.RegisteredClaims
}

// GetNewToken
// generate new token
func GetNewToken(cred Credential) (token string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Credential: cred,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenJWT.SignedString(secretKey)
	return
}

func ReadAndVerifyToken(token string) (Claims, error) {
	var claims Claims
	tokenJWT, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !tokenJWT.Valid {
		return Claims{}, err
	}

	return claims, nil
}
