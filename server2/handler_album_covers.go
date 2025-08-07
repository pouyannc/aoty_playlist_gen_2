package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/pouyannc/aoty_list_gen/internal/scrape"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
)

type AlbumCoversResp struct {
	ID       string `json:"id"`
	Artist   string `json:"artist"`
	ImageURL string `json:"image_url"`
}

func (cfg *apiConfig) handlerAlbumCovers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	nrCovers := 6

	type scrapeParams struct {
		scrapeURL string
		filter    string
	}

	query := r.URL.Query()
	params := scrapeParams{
		scrapeURL: query.Get("scrape_url"),
		filter:    query.Get("type"),
	}

	allScrapeURLs, err := scrape.CreateAllScrapeURLs(params.scrapeURL, params.filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create scrape urls", err)
		return
	}

	fmt.Println("Filter: ", params.filter)

	var pages []*rod.Page

	for _, url := range allScrapeURLs {
		page := stealth.MustPage(cfg.browser)
		defer page.MustClose()
		page.MustNavigate(url)
		fmt.Println(url)
		pages = append(pages, page)
	}
	fmt.Println("========== Opened pages in:", time.Since(startTime))
	startTime = time.Now()

	albums, err := scrape.ScrapeAlbums(pages, params.filter, nrCovers)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't scrape albums", err)
		return
	}
	fmt.Printf("========== Scraped in: %v\n", time.Since(startTime))
	startTime = time.Now()

	session, err := cfg.store.Get(r, "spotify-session")
	if err != nil {
		fmt.Println(err)
		return
	}
	token, ok := session.Values["access_token"]
	if !ok {
		fmt.Println(token)
		return
	}

	albumData, err := spotify.AlbumData(albums, token.(string))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("========== Spotify search done in: %v\n", time.Since(startTime))

	resp := []AlbumCoversResp{}
	for _, data := range albumData {
		resp = append(resp, AlbumCoversResp{
			ID:       data.AlbumID,
			Artist:   data.Artist,
			ImageURL: data.CoverURL,
		})

		if len(resp) == nrCovers {
			break
		}
	}

	respondWithJSON(w, http.StatusOK, resp)
}
