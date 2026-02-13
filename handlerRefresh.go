package main

import (
	"main/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type Token struct {
		Token string `json:"token"`
	}

	tokenBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenBearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Token Expired")
		return
	}

	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token Revoked")
		return
	}

	newToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := Token{
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
