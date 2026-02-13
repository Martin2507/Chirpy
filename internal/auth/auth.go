package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		return "", err
	}

	return hashed_password, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	boolResult, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false, err
	}

	return boolResult, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	registeredClaim := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaim)

	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy-access" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Unauthorized Access")
	}

	splitToken := strings.Split(authHeader, " ")

	return splitToken[1], nil

}

func MakeRefreshToken() (string, error) {

	data := make([]byte, 32)

	_, err := rand.Read(data)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Unauthorized Access")
	}

	splitPolkaKey := strings.Split(authHeader, " ")

	if len(splitPolkaKey) < 2 {
		return "", fmt.Errorf("Unauthorized Access")
	}

	if splitPolkaKey[0] != "ApiKey" {
		return "", fmt.Errorf("Unauthorized Access")
	}

	return splitPolkaKey[1], nil
}
