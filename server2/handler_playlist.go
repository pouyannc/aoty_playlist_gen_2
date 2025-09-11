package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pouyannc/aoty_list_gen/internal/scrape"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
	"github.com/pouyannc/aoty_list_gen/util"
)

type PlaylistData struct {
	PlaylistID string `json:"playlist_id"`
}

type scrapeParamsPlaylist struct {
	scrapeURL      string
	nTracks        int
	tracksPerAlbum int
	filter         string
}

func (cfg *apiConfig) handlerPlaylist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rParams struct {
		PlaylistName string `json:"playlistName"`
	}
	err := decoder.Decode(&rParams)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	query := r.URL.Query()
	nTracksInt, err := strconv.Atoi(query.Get("nr_tracks"))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse number of tracks to int", err)
		return
	}
	tracksPerInt, err := strconv.Atoi(query.Get("tracks_per"))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse tracks per to int", err)
		return
	}
	qParams := scrapeParamsPlaylist{
		scrapeURL:      query.Get("scrape_url"),
		nTracks:        nTracksInt,
		tracksPerAlbum: tracksPerInt,
		filter:         query.Get("type"),
	}

	fmt.Println("REQUEST QUERY PARAMS:", qParams)

	nAlbums := (qParams.nTracks / qParams.tracksPerAlbum) + 1

	allScrapeURLs, err := scrape.CreateAllScrapeURLs(qParams.scrapeURL, qParams.filter)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to create scrape urls", err)
		return
	}

	page := cfg.browser.MustPage("https://www.albumoftheyear.org/")
	defer page.MustClose()

	albums, _ := scrape.ScrapeAlbums(page, allScrapeURLs, qParams.filter, nAlbums)

	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Unable to get server session", err)
		return
	}
	token, ok := session.Values["access_token"]
	if !ok {
		util.RespondWithError(w, http.StatusUnauthorized, "No access token found in user session", fmt.Errorf("no access token in session: %v", token))
		return
	}
	uid, ok := session.Values["spotify_uid"]
	if !ok {
		util.RespondWithError(w, http.StatusUnauthorized, "No spotify uid found in user session", fmt.Errorf("no uid in session: %v", uid))
		return
	}

	albumData, err := spotify.AlbumData(albums, token.(string))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get album data from Spotify", err)
		return
	}

	trackURIs, err := spotify.GetTracklist(albumData, tracksPerInt, nTracksInt, token.(string))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't get track URIs from Spotify", err)
		return
	}

	playlistID, err := spotify.CreatePlaylist(trackURIs, token.(string), uid.(string), rParams.PlaylistName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't create or populate playlist", err)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, PlaylistData{
		PlaylistID: playlistID,
	})
}
