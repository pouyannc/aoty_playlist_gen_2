package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pouyannc/aoty_list_gen/util"
)

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
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
				token, err := refreshAndGetTokens(refreshToken.(string))
				if err != nil {
					util.RespondWithError(w, http.StatusInternalServerError, "Couldn't refresh tokens", err)
					return
				}

				session.Values["access_token"] = token.AccessToken
				session.Values["refresh_token"] = token.RefreshToken
				session.Values["expiry"] = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
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
	url := "https://accounts.spotify.com/api/token"

	payload := struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_type"`
	}{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return TokenResp{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return TokenResp{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
