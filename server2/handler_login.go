package main

import (
	"net/http"

	"github.com/pouyannc/aoty_list_gen/util"
	"golang.org/x/oauth2"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	state, err := util.GenerateRandomString(16)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to generate state for Spotify OAuth2", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Path:     "/api/login",
	})

	url := cfg.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
