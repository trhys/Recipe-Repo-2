package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

func SeedDB(file, ik string, db *sql.DB, ctx context.Context) error {
	log.Println("Performing initial setup:")
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Panic("Failed to load init file")
	}

	log.Println("File read. Marshalling...")
	var ings struct {
		Ingredients []struct{
			Name string `json:"name"`
		} `json:"ingredients"`
	}

	if err := json.Unmarshal(bytes, &ings); err != nil {
		log.Panic("Failed to read file into struct")
	}

	log.Println("Successfully marshalled json: creating ingredients in database...")

	dbConn := database.New(db)

	for _, i := range ings.Ingredients {
		query := database.CreateIngredientParams{
			Name: i.Name,
			ImageKey: ik,
		}

		if _, err := dbConn.CreateIngredient(ctx, query); err != nil {
			log.Panic("Couldnt create ingredient during setup")
		}
	}

	log.Println("Setup successful...")
	return nil
}
