package main

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/martin/Chirpy/internal/database"
)

type apiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	Secret         string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type UserParams struct {
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
}

type Payload struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
