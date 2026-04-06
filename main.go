package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"database/sql"

	_ "github.com/lib/pq"
        "github.com/joho/godotenv"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB")
	if dbUrl == "" {
		log.Fatal("Failed to load database: url missing")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Failed to load platform config")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("Failed to load secret")
	}

	jwtDur := os.Getenv("JWT_DUR")
	if jwtDur == "" {
		log.Fatal("Failed to load jwt duration")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to load database: connection failed")
	}

	convDur, err := strconv.Atoi(jwtDur)
	if err != nil {
		log.Print("Failed to load jwt duration - defaulting to 3600")
		convDur = 3600
	}
	jwtDuration := time.Duration(convDur)*time.Second

	// Load api config
	
	config := apiConfig{
		db: database.New(db),
		platform: platform,
		secret: secret,
		jwtDuration: jwtDuration,
	}

	log.Print("Successfully loaded database...")

	// Load server

	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	// Handlers

	mux.HandleFunc("GET /api/users/{user_id}", config.handlerGetUser)
	mux.HandleFunc("POST /api/new_user", config.handlerCreateUser)
	mux.HandleFunc("POST /api/login", config.handlerLogin)
	mux.HandleFunc("POST /api/reset", config.handlerReset)

	mux.HandleFunc("GET /api/recipes/{recipe_id}", config.handlerGetRecipe)
	mux.HandleFunc("POST /api/new_recipe", config.handlerCreateRecipe)

	log.Print("Successfully loaded server config...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
