package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CheckPasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func HashPassword(passowrd string) (string, error) {
	return argon2id.CreateHash(passowrd, argon2id.DefaultParams)
}

func MakeJWT(userID uuid.UUID, tokenString string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(tokenString))
	if err != nil {
		return "", err
	}

	return t, nil

}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	returnData, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return returnData, nil

}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Unauthorized Access")
	}

	splitToken := strings.Split(authHeader, " ")

	return splitToken[1], nil

}
