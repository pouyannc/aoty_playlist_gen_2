package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/pouyannc/aoty_list_gen/internal/spotify"
)

func (cfg *apiConfig) handlerLoginCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Spotify Oauth callback code", errors.New("missing code"))
		return
	}

	token, err := cfg.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Spotify Oauth token exchange failed", err)
		return
	}

	uid, err := spotify.GetUID(token.AccessToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get spotify UID", err)
		return
	}

	session, _ := cfg.store.Get(r, "spotify-session")
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expires_in"] = token.ExpiresIn
	session.Values["spotify_uid"] = uid
	err = session.Save(r, w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save server session", err)
		return
	}

	http.Redirect(w, r, "http://localhost:5173/auth/callback", http.StatusFound)
}
