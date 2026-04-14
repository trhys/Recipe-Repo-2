package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"database/sql"

	_ "github.com/lib/pq"
        "github.com/joho/godotenv"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
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

	convDur, err := strconv.Atoi(jwtDur)
	if err != nil {
		log.Print("Failed to load jwt duration - defaulting to 3600")
		convDur = 3600
	}
	jwtDuration := time.Duration(convDur)*time.Second


	appDirectory := os.Getenv("APP_DIR")
	if appDirectory == "" {
		log.Fatal("Failed to load app directory")
	}
	
	s3bucket := os.Getenv("S3_BUCKET")
	if s3bucket == "" {
		log.Fatal("Failed to load s3 bucket")
	}

	s3region := os.Getenv("S3_REGION")
	if err != nil {
		log.Fatal("Failed to load s3 region")
	}

	s3cdn := os.Getenv("S3_CDN")
	if s3cdn == "" {
		log.Fatal("Failed to load s3 CDN")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to load database: connection failed")
	}

	log.Print("Successfully loaded database...")

	// Load S3 client
	s3cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s3region))
	if err != nil {
		log.Fatal("Failed to load s3 config")
	}

	log.Print("Successfully loaded s3 config...")

	// Load api config
	
	config := apiConfig{
		db: database.New(db),
		platform: platform,
		secret: secret,
		jwtDuration: jwtDuration,
		s3client: s3.NewFromConfig(s3cfg) ,
		s3bucket: s3bucket,
		s3region: s3region,
		s3cdn: s3cdn,
	}
	
	// Load server

	mux := http.NewServeMux()
	server := http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
	}

	// JS Fileserver handler
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(appDirectory)))
	mux.Handle("/app/", appHandler)

	// Handlers
	mux.HandleFunc("GET /app/recipes/{recipe_id}", config.appGetRecipe)
	mux.HandleFunc("GET /app/recipes/by_user/{user_id}", config.appGetUsersRecipes)

	mux.HandleFunc("GET /api/users/{user_id}", config.handlerGetUser)
	mux.HandleFunc("POST /api/new_user", config.handlerCreateUser)
	mux.HandleFunc("POST /api/login", config.handlerLogin)
	mux.HandleFunc("POST /api/reset", config.handlerReset)

	mux.HandleFunc("GET /api/recipes/by_user/{user_id}", config.handlerGetUsersRecipes)
	mux.HandleFunc("GET /api/recipes/{recipe_id}", config.handlerGetRecipe)
	mux.HandleFunc("GET /api/recipes", config.handlerGetRecipeList)
	mux.HandleFunc("POST /api/new_recipe", config.handlerCreateRecipe)

	log.Print("Successfully loaded server...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
