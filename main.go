package main

import (
	"log"
	"net/http"
	"os"
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

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to load database: connection failed")
	}

	// Load api config
	
	config := apiConfig{
		db: database.New(db),
	}

	log.Print("Successfully loaded database...")

	// Load server

	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	// Handlers

	mux.HandleFunc("POST /api/new_user", config.handlerCreateUser)

	log.Print("Successfully loaded server config...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
