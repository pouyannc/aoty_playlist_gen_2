package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/pouyannc/aoty_list_gen/internal/middleware"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type apiConfig struct {
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	browser     *rod.Browser
	rdb         *redis.Client
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func main() {
	gob.Register(time.Time{})

	_ = godotenv.Load()
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	log.Fatal("error loading .env (port is empty)")
	// }

	redisAddr := os.Getenv("REDIS_ADDR")
	fmt.Println("redis address:", redisAddr)
	//redisAddr := ":6379" //for dev

	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		log.Fatalf("Invalid REDIS_ADDR: %v", err)
	}
	rdb := redis.NewClient(opt)
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
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	u := launcher.
		NewUserMode().
		Headless(true).
		Leakless(true).
		UserDataDir("tmp/t").
		Set("disable-default-apps").
		Set("no-first-run").
		Set("disable-gpu").
		NoSandbox(true).
		MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	cfg := apiConfig{
		oauthConfig: oauthConfig,
		store:       store,
		browser:     browser,
		rdb:         rdb,
	}

	spa := spaHandler{
		staticPath: "dist",
		indexPath:  "index.html",
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/login", cfg.handlerLogin).Methods("GET")
	r.HandleFunc("/api/login/callback", cfg.handlerLoginCallback).Methods("GET")
	r.HandleFunc("/api/login/guest", cfg.handlerLoginGuest).Methods("GET")
	r.HandleFunc("/api/auth/tokens", cfg.handlerAuthTokens).Methods("GET")
	r.HandleFunc("/api/logout", cfg.handlerLogout).Methods("DELETE")

	albumsSubrouter := r.PathPrefix("/api/albums").Subrouter()
	albumsSubrouter.HandleFunc("/covers", cfg.handlerAlbumCovers).Methods("GET")
	albumsSubrouter.HandleFunc("/playlist", cfg.handlerPlaylist).Methods("POST")
	albumsSubrouter.Use(middleware.ValidateSpotifyToken(cfg.store))

	r.PathPrefix("/").Handler(spa)

	// corsHandler := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{os.Getenv("CLIENT_URL")},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// }).Handler(r)
	// Old CORS code, used for SPA being run on separate server

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Printf("Server running: port %v", server.Addr)
	log.Fatal(server.ListenAndServe())
}
