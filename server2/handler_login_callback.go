package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/pouyannc/aoty_list_gen/internal/spotify"
	"github.com/pouyannc/aoty_list_gen/util"
)

func (cfg *apiConfig) handlerLoginCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing Spotify Oauth callback code", errors.New("missing code"))
		return
	}

	token, err := cfg.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Spotify Oauth token exchange failed", err)
		return
	}

	uid, err := spotify.GetUID(token.AccessToken)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get spotify UID", err)
		return
	}

	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get or create server spotify-session", err)
		return
	}
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	session.Values["spotify_uid"] = uid
	err = session.Save(r, w)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't save server session", err)
		return
	}

	clientURL := os.Getenv("CLIENT_URL")
	http.Redirect(w, r, clientURL, http.StatusFound)
}
