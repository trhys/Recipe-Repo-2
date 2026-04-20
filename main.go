package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
        "github.com/joho/godotenv"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	"github.com/trhys/Recipe-Repo-2/internal/migrations"
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

	adminDir := os.Getenv("ADMIN_DIR")
	if adminDir == "" {
		log.Fatal("Failed to load admin directory")
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

	imagePlaceholder := os.Getenv("IMAGE_PLACEHOLDER")
	if imagePlaceholder == "" {
		log.Fatal("Failed to load placehold for images")
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
	
	cfg := apiConfig{
		db: database.New(db),
		platform: platform,
		secret: secret,
		jwtDuration: jwtDuration,
		s3client: s3.NewFromConfig(s3cfg) ,
		s3bucket: s3bucket,
		s3region: s3region,
		s3cdn: s3cdn,
		imagePlaceholder: imagePlaceholder,
	}

	// Check database seeding
	log.Println("Checking database seeding...")
	var count int
	db.QueryRow("SELECT COUNT(*) FROM ingredients").Scan(&count)
	if count < 99 {
		log.Println("Unseeded!")
		if err := migrations.SeedDB("setup.json", imagePlaceholder, db, context.Background()); err != nil {
			log.Panic("Failed to seed database")
		}
	} else {
		log.Println("Database already seeded, continue...")
	}
	
	// Load server

	mux := http.NewServeMux()
	server := http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
	}

	// JS Fileserver handler
	appHandler := http.FileServer(http.Dir(appDirectory))
	mux.Handle("/", appHandler)

	adminHandler := http.StripPrefix("/admin", http.FileServer(http.Dir(adminDir)))
	mux.Handle("/admin/", adminHandler)

	// Handlers :

	// Web app eps
	mux.HandleFunc("GET /recipes/{recipe_id}", cfg.handlerGetRecipe)
	mux.HandleFunc("GET /users/{user_id}", cfg.handlerGetUsersRecipes)
	mux.HandleFunc("GET /shoppinglists/{shopping_list_id}", cfg.handlerGetShoppingList)

	// User eps
	mux.HandleFunc("GET /api/users/{user_id}", cfg.handlerGetUser)
	mux.HandleFunc("POST /api/new_user", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/admin/reset", cfg.handlerReset)

	// Recipe eps
	mux.HandleFunc("GET /api/recipes/by_user/{user_id}", cfg.handlerGetUsersRecipes)
	mux.HandleFunc("GET /api/recipes/{recipe_id}", cfg.handlerGetRecipe)
	mux.HandleFunc("GET /api/recipes", cfg.handlerGetRecipeList)
	mux.HandleFunc("POST /api/new_recipe", cfg.handlerCreateRecipe)

	// Ingredient eps
	mux.HandleFunc("POST /api/admin/new_ingredient", cfg.handlerCreateIngredient)
	mux.HandleFunc("GET /api/get_ingredients", cfg.handlerGetIngredientBase)

	// Shopping list eps
	mux.HandleFunc("POST /api/new_shopping_list", cfg.handlerCreateShoppingList)
	mux.HandleFunc("POST /api/add_to_list", cfg.handlerAddToShoppingList)
	mux.HandleFunc("GET /api/shoppinglists/{shopping_list_id}", cfg.handlerGetShoppingList)

	// Token eps
	mux.HandleFunc("POST /api/tokens/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/tokens/revoke", cfg.handlerRevokeToken)

	// :

	log.Print("Successfully loaded server...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
