package main

import (
	"fmt"
	"net/http"
	"net/url"
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

	fmt.Println(params.filter)
	nPages := 1
	switch params.filter {
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
		page.MustNavigate(params.scrapeURL)
		pages[i] = page

		if nPages > 1 {
			err := page.Timeout(5*time.Second).WaitElementsMoreThan(".prev", 0)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Failed to load albums page", err)
				return
			}

			baseURL, _ := url.Parse(page.MustInfo().URL)
			fullUrl, _ := baseURL.Parse(*page.MustElement(`a:has(div.prev)`).MustAttribute("href"))
			params.scrapeURL = string(fullUrl.String())
		}
	}
	fmt.Printf("Opened pages in: %v", time.Since(startTime))

	albums, _ := scrape.ScrapeAlbums(pages, params.filter, nrCovers)
	fmt.Printf("Scraped in: %v", time.Since(startTime))
	// for debugging...
	fmt.Println(albums)
	fmt.Println(params)

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

	fmt.Printf("Spotify search done in: %v", time.Since(startTime))

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
