package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Arrange
	secret := "my-secret-key"
	userID := uuid.New()
	expiresIn := time.Hour

	// Act
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	// Assert
	if validatedID != userID {
		t.Errorf("expected %v, got %v", userID, validatedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "my-secret-key"
	userID := uuid.New()

	// Create a token that was "born" an hour ago and is already expired
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	// We EXPECT an error here!
	if err == nil {
		t.Errorf("expected error for expired token, but got nil")
	}
}
