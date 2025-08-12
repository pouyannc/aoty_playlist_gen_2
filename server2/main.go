package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type apiConfig struct {
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	browser     *rod.Browser
	rdb         *redis.Client
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis at %s: %v", redisAddr, err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"playlist-modify-public"},
		Endpoint:     spotify.Endpoint,
	}

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	u := launcher.
		NewUserMode().
		Headless(true).
		Leakless(true).
		UserDataDir("tmp/t").
		Set("disable-default-apps").
		Set("no-first-run").
		MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	cfg := apiConfig{
		oauthConfig: oauthConfig,
		store:       store,
		browser:     browser,
		rdb:         rdb,
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/login/refresh", cfg.handlerRefreshLogin).Methods("GET")
	r.HandleFunc("/api/login", cfg.handlerLogin).Methods("GET")
	r.HandleFunc("/api/login/callback", cfg.handlerLoginCallback).Methods("GET")
	r.HandleFunc("/api/auth/tokens", cfg.handlerAuthTokens).Methods("GET")

	r.HandleFunc("/api/albums/covers", cfg.handlerAlbumCovers).Methods("GET")
	r.HandleFunc("/api/albums/playlist", cfg.handlerPlaylist).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsHandler,
	}

	log.Printf("Server running: http://localhost%v", server.Addr)
	log.Fatal(server.ListenAndServe())
}
