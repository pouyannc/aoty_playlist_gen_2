package main

import (
	"net/http"
	"os"
	"time"

	"github.com/pouyannc/aoty_list_gen/internal/middleware"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
	"github.com/pouyannc/aoty_list_gen/util"
)

func (cfg *apiConfig) handlerLoginGuest(w http.ResponseWriter, r *http.Request) {
	guestRefreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")
	tokens, err := middleware.RefreshAndGetTokens(guestRefreshToken)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get auth tokens", err)
		return
	}

	uid, err := spotify.GetUID(tokens.AccessToken)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get spotify UID", err)
		return
	}
	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get or create server spotify-session", err)
		return
	}
	session.Values["access_token"] = tokens.AccessToken
	session.Values["refresh_token"] = guestRefreshToken
	session.Values["expiry"] = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	session.Values["spotify_uid"] = uid
	err = session.Save(r, w)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't save server session", err)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, struct{}{})
}
