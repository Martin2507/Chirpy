package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "well-kept-secret"
	userID := uuid.New()

	// 1. Test Valid JWT
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("failed to validate valid JWT: %v", err)
	}
	if validatedID != userID {
		t.Errorf("expected userID %v, got %v", userID, validatedID)
	}

	// 2. Test Wrong Secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("expected error when validating with wrong secret, got nil")
	}

	// 3. Test Expired Token
	// We create a token that expired 1 hour ago
	expiredToken, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to make expired JWT: %v", err)
	}

	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("expected error when validating expired token, got nil")
	}
}
