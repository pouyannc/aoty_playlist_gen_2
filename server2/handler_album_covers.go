package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/gorilla/sessions"
	"github.com/pouyannc/aoty_list_gen/internal/scrape"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
	"github.com/redis/go-redis/v9"
)

type AlbumCoversResp struct {
	ID       string `json:"id"`
	Artist   string `json:"artist"`
	ImageURL string `json:"image_url"`
}

type cachePayload struct {
	Data []AlbumCoversResp `json:"data"`
	Ts   int64             `json:"ts"`
}

type fetchAndCacheError struct {
	err    error
	status int
}

var (
	cacheKey  = "albumCovers"
	freshness = 4 * time.Hour
)

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

	var resp []AlbumCoversResp
	var fetchAndCacheErr fetchAndCacheError

	key := fmt.Sprintf("%s:%s:%s", cacheKey, params.scrapeURL, params.filter)
	cacheValue, err := cfg.rdb.Get(context.Background(), key).Result()
	if err == nil {
		var payload cachePayload
		err := json.Unmarshal([]byte(cacheValue), &payload)
		if err == nil {
			respondWithJSON(w, http.StatusOK, payload.Data)
			payloadAge := time.Since(time.Unix(payload.Ts, 0))
			if payloadAge > freshness {
				go fetchAndCacheAlbumData(*r, cfg.browser, cfg.store, cfg.rdb, params.scrapeURL, params.filter, key, &resp, &fetchAndCacheErr)
			}
			return
		}
	}

	fetchAndCacheAlbumData(*r, cfg.browser, cfg.store, cfg.rdb, params.scrapeURL, params.filter, key, &resp, &fetchAndCacheErr)

	if fetchAndCacheErr.err != nil {
		respondWithError(w, fetchAndCacheErr.status, "Couldn't fetch or cache album cover data", fetchAndCacheErr.err)
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func fetchAndCacheAlbumData(r http.Request, browser *rod.Browser, store *sessions.CookieStore, rdb *redis.Client, scrapeURL, filter, key string, resp *[]AlbumCoversResp, fcErr *fetchAndCacheError) {
	startTime := time.Now()

	nrCovers := 6

	allScrapeURLs, err := scrape.CreateAllScrapeURLs(scrapeURL, filter)
	if err != nil {
		fcErr.status, fcErr.err = http.StatusInternalServerError, err
		return
	}

	fmt.Println("Filter: ", filter)

	var pages []*rod.Page

	for _, url := range allScrapeURLs {
		page := stealth.MustPage(browser)
		defer page.MustClose()
		page.MustNavigate(url)
		fmt.Println(url)
		pages = append(pages, page)
	}
	fmt.Println("========== Opened pages in:", time.Since(startTime))
	startTime = time.Now()

	albums, err := scrape.ScrapeAlbums(pages, filter, nrCovers)
	if err != nil {
		fcErr.status, fcErr.err = http.StatusInternalServerError, err
		return
	}
	fmt.Printf("========== Scraped in: %v\n", time.Since(startTime))
	startTime = time.Now()

	session, err := store.Get(&r, "spotify-session")
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

	for _, data := range albumData {
		*resp = append(*resp, AlbumCoversResp{
			ID:       data.AlbumID,
			Artist:   data.Artist,
			ImageURL: data.CoverURL,
		})

		if len(*resp) == nrCovers {
			break
		}
	}

	payload := cachePayload{
		Data: *resp,
		Ts:   time.Now().Unix(),
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error marshaling cache payload: %v\n", err)
	}
	err = rdb.Set(context.Background(), key, bytes, 0).Err()
	if err != nil {
		fmt.Printf("error saving response to redis cache: %v\n", err)
	}
}
