package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	
	"github.com/trhys/Recipe-Repo-2/internal/database"
	pb "github.com/schollz/progressbar/v3"
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
			Conversions []struct {
				From string `json:"from_unit"`
				To string `json:"to_unit"`
				Ratio float32 `json:"ratio"`
			} `json:"conversions"`
		} `json:"ingredients"`
	}

	if err := json.Unmarshal(bytes, &ings); err != nil {
		log.Panic("Failed to read file into struct")
	}

	log.Println("Successfully marshalled json: creating ingredients in database...")

	dbConn := database.New(db)

	bar := pb.Default(int64(len(ings.Ingredients)))
	for _, i := range ings.Ingredients {
		queryA := database.CreateIngredientParams{
			Name: i.Name,
			ImageKey: ik,
		}

		ingredient, err := dbConn.CreateIngredient(ctx, queryA) 
		if err != nil {
			log.Printf("Failed on ingredient: %s - ERROR: %v", i.Name, err)
			log.Panic("Couldnt create ingredient during setup")
		}

		for index, conv := range i.Conversions {
			queryB := database.CreateConversionParams{
				IngredientID: ingredient.ID,
				FromUnit: conv.From,
				ToUnit: conv.To,
				Ratio: conv.Ratio,
			}

			if err := dbConn.CreateConversion(ctx, queryB); err != nil {
				log.Printf("Failed to create conversion at position: %d - %s", index, i.Name)
			}
		}
		bar.Add(1)
	}

	log.Println("Setup successful...")
	return nil
}
