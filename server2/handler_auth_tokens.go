package main

import (
	"errors"
	"net/http"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int64 `json:"expires_in"`
}

func (cfg *apiConfig) handlerAuthTokens(w http.ResponseWriter, r *http.Request) {
	session, _  := cfg.store.Get(r, "spotify-session")
	access, ok1 := session.Values["access_token"].(string)
	refresh, ok2 := session.Values["refresh_token"].(string)
	expires, ok3 := session.Values["expires_in"].(int64)

	if !ok1 || !ok2 || !ok3 {
		respondWithError(w, http.StatusUnauthorized, "User unauthorized", errors.New("tokens not found in session"))
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		AccessToken: access,
		RefreshToken: refresh,
		ExpiresIn: expires,
	})
}