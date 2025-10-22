package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pouyannc/aoty_list_gen/util"
)

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	SpotifyUID   string    `json:"spotify_uid"`
}

func (cfg *apiConfig) handlerAuthTokens(w http.ResponseWriter, r *http.Request) {
	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Couldn't decode existing spotify session", err)
		return
	}
	access, ok1 := session.Values["access_token"].(string)
	refresh, ok2 := session.Values["refresh_token"].(string)
	expiry, ok3 := session.Values["expiry"].(time.Time)
	uid, ok4 := session.Values["spotify_uid"].(string)

	fmt.Println("Cookie data:", access, refresh, expiry, uid)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		util.RespondWithError(w, http.StatusUnauthorized, "User unauthorized", errors.New("tokens not found in session"))
		return
	}

	util.RespondWithJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		Expiry:       expiry,
		SpotifyUID:   uid,
	})
}
