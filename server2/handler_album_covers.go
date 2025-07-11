package main

import (
	"fmt"
	"net/http"

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
	type scrapeParams struct {
		scrapeURL string
		filter    string
	}

	query := r.URL.Query()
	params := scrapeParams{
		scrapeURL: query.Get("scrape_url"),
		filter:    query.Get("type"),
	}

	fmt.Println("---- opening new page...")
	page := stealth.MustPage(cfg.browser)
	defer page.MustClose()
	page.MustNavigate(params.scrapeURL)

	albums, _ := scrape.ScrapeAlbums(page, 6)

	// for debugging...
	fmt.Println(albums)
	fmt.Println(params)
	page.MustScreenshot("debug.png")

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

	fmt.Println(albumData)

	// marshall and pass through json response
	resp := []AlbumCoversResp{}
	for _, data := range albumData {
		resp = append(resp, AlbumCoversResp{
			ID:       data.AlbumID,
			Artist:   data.Artist,
			ImageURL: data.CoverURL,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}
