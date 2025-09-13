package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pouyannc/aoty_list_gen/util"
)

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type contextKey string

const TokenKey contextKey = "token"

func ValidateSpotifyToken(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "spotify-session")
			if err != nil {
				util.RespondWithError(w, http.StatusBadRequest, "Failed to find spotify session store", err)
				return
			}

			expiry, ok := session.Values["expiry"]
			if !ok {
				util.RespondWithError(w, http.StatusBadRequest, "No token expiry found in spotify session store", errors.New("expiry key not found in spotify session store"))
				return
			}

			refreshToken, ok := session.Values["refresh_token"]
			if !ok {
				util.RespondWithError(w, http.StatusBadRequest, "No refresh token found in spotify session store", errors.New("refresh token not found in spotify session store"))
				return
			}
			if time.Now().After(expiry.(time.Time)) {
				tokens, err := refreshAndGetTokens(refreshToken.(string))
				if err != nil {
					util.RespondWithError(w, http.StatusInternalServerError, "Couldn't refresh tokens", err)
					return
				}

				session.Values["access_token"] = tokens.AccessToken
				session.Values["expiry"] = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
				err = session.Save(r, w)
				if err != nil {
					util.RespondWithError(w, http.StatusInternalServerError, "Couldn't save server session", err)
					return
				}
			}

			accessToken, ok := session.Values["access_token"].(string)
			if !ok {
				util.RespondWithError(w, http.StatusBadRequest, "No access token found in spotify session store", errors.New("access token not found in spotify session store"))
				return
			}

			ctx := context.WithValue(r.Context(), TokenKey, accessToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func refreshAndGetTokens(refreshToken string) (TokenResp, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	formData := url.Values{}
	formData.Add("grant_type", "refresh_token")
	formData.Add("refresh_token", refreshToken)
	encodedData := formData.Encode()

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(encodedData))
	if err != nil {
		return TokenResp{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	authClientStr := fmt.Sprintf("%s:%s", clientID, clientSecret)
	authClientEnc := base64.URLEncoding.EncodeToString([]byte(authClientStr))
	authHeader := fmt.Sprintf("Basic %s", authClientEnc)
	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return TokenResp{}, fmt.Errorf("spotify returned status %d, and error reading body: %v", resp.StatusCode, err)
		}
		return TokenResp{}, fmt.Errorf("spotify returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	decoder := json.NewDecoder(resp.Body)
	var token TokenResp
	err = decoder.Decode(&token)
	if err != nil {
		return TokenResp{}, err
	}

	return token, nil
}
