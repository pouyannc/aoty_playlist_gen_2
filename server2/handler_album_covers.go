package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pouyannc/aoty_list_gen/internal/middleware"
	"github.com/pouyannc/aoty_list_gen/internal/spotify"
	"github.com/pouyannc/aoty_list_gen/util"
)

type AlbumCoversResp struct {
	ID       string `json:"id"`
	Artist   string `json:"artist"`
	ImageURL string `json:"image_url"`
}

func (cfg *apiConfig) handlerAlbumCovers(w http.ResponseWriter, r *http.Request) {
	nCovers := 8

	query := r.URL.Query()
	albumScrapeKeyParam := query.Get("scrape_key")

	var resp []AlbumCoversResp

	key := fmt.Sprintf("%s:%s", cacheScrapeKey, albumScrapeKeyParam)
	cacheValue, err := cfg.rdb.Get(context.Background(), key).Result()
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve redis album data", err)
		return
	}

	var payload cacheScrapePayload
	err = json.Unmarshal([]byte(cacheValue), &payload)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Couldn't unmarshal cache scrape payload", err)
	}

	token := r.Context().Value(middleware.TokenKey).(string)

	albumData, err := spotify.AlbumData(payload.ScrapeAlbums, token, nCovers)
	if err != nil {
		fmt.Println(err)
	}

	for _, data := range albumData {
		resp = append(resp, AlbumCoversResp{
			ID:       data.AlbumID,
			Artist:   data.Artist,
			ImageURL: data.CoverURL,
		})

		if len(resp) == nCovers {
			break
		}
	}

	// payload := cachePayload{
	// 	Data: *resp,
	// 	Ts:   time.Now().Unix(),
	// }
	// bytes, err := json.Marshal(payload)
	// if err != nil {
	// 	fmt.Printf("error marshaling cache payload: %v\n", err)
	// }
	// err = rdb.Set(context.Background(), key, bytes, 0).Err()
	// if err != nil {
	// 	fmt.Printf("error saving response to redis cache: %v\n", err)
	// }

	fmt.Println("Album cover data retrieved -------------------------------------------------------------------")

	util.RespondWithJSON(w, http.StatusOK, resp)
}
