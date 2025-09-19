package main

import (
	"net/http"

	"github.com/pouyannc/aoty_list_gen/util"
)

func (cfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request) {
	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get the server session", err)
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't save the updated server session", err)
		return
	}
}
