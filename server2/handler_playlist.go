package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/pouyannc/aoty_list_gen/internal/scrape"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
)

type PlaylistData struct {
	PlaylistID string `json:"playlist_id"`
}

func (cfg *apiConfig) handlerPlaylist(w http.ResponseWriter, r *http.Request) {
	type scrapeParams struct {
		scrapeURL      string
		nTracks        int
		tracksPerAlbum int
		filter         string
	}

	decoder := json.NewDecoder(r.Body)
	var rParams struct {
		PlaylistName string `json:"playlistName"`
	}
	err := decoder.Decode(&rParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	query := r.URL.Query()
	nTracksInt, err := strconv.Atoi(query.Get("nr_tracks"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse number of tracks to int", err)
		return
	}
	tracksPerInt, err := strconv.Atoi(query.Get("tracks_per"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse tracks per to int", err)
		return
	}
	qParams := scrapeParams{
		scrapeURL:      query.Get("scrape_url"),
		nTracks:        nTracksInt,
		tracksPerAlbum: tracksPerInt,
		filter:         query.Get("type"),
	}

	fmt.Println("REQUEST QUERY PARAMS:", qParams)

	nAlbums := (qParams.nTracks / qParams.tracksPerAlbum) + 1

	nPages := 1
	switch qParams.filter {
	case "months":
		nPages = 4
	case "years":
		nPages = 3
	}

	pages := make([]*rod.Page, nPages)

	for i := 0; i < nPages; i++ {
		fmt.Println("---- opening new page...")
		page := stealth.MustPage(cfg.browser)
		defer page.MustClose()
		page.MustNavigate(qParams.scrapeURL)
		pages[i] = page

		if nPages > 1 {
			err := page.Timeout(5*time.Second).WaitElementsMoreThan(".prev", 0)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to load albums page", err)
				return
			}

			baseURL, _ := url.Parse(page.MustInfo().URL)
			fullUrl, _ := baseURL.Parse(*page.MustElement(`a:has(div.prev)`).MustAttribute("href"))
			qParams.scrapeURL = string(fullUrl.String())
		}
	}

	albums, _ := scrape.ScrapeAlbums(pages, qParams.filter, nAlbums)

	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get server session", err)
		return
	}
	token, ok := session.Values["access_token"]
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "No access token found in user session", fmt.Errorf("no access token in session: %v", token))
		return
	}
	uid, ok := session.Values["spotify_uid"]
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "No spotify uid found in user session", fmt.Errorf("no uid in session: %v", uid))
		return
	}

	albumData, err := spotify.AlbumData(albums, token.(string))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get album data from Spotify", err)
		return
	}

	trackURIs, err := spotify.GetTracklist(albumData, tracksPerInt, nTracksInt, token.(string))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get track URIs from Spotify", err)
		return
	}

	playlistID, err := spotify.CreatePlaylist(trackURIs, token.(string), uid.(string), rParams.PlaylistName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create or populate playlist", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, PlaylistData{
		PlaylistID: playlistID,
	})
}
