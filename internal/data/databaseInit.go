package data

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"
	
	"github.com/lib/pq"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	pb "github.com/schollz/progressbar/v3"
)

//go:embed setup.json
var setup []byte

func InitDB(ik string, db *sql.DB, ctx context.Context) error {
	log.Println("Checking database...")

	log.Println("Loading ingredients from JSON...")
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

	if err := json.Unmarshal(setup, &ings); err != nil {
		log.Panic("Failed to marshal JSON!")
	}

	log.Println("Successfully read file - verifying entries...")

	dbConn := database.New(db)

	bar := pb.Default(int64(len(ings.Ingredients)))
	for _, i := range ings.Ingredients {
		ingID, err := dbConn.GetIngredientFromName(ctx, i.Name)
		if err == nil {
			c, err := dbConn.GetConversionsByID(ctx, ingID)
			if err == nil && len(c) == len(i.Conversions) {
				bar.Add(1)
				continue
			} else {
				for index, conv := range i.Conversions {
					queryB := database.CreateConversionParams{
						IngredientID: ingID,
						FromUnit: conv.From,
						ToUnit: conv.To,
						Ratio: conv.Ratio,
					}
				
					if err := dbConn.CreateConversion(ctx, queryB); err != nil {
						if pqErr, ok := err.(*pq.Error); ok {
							if pqErr.Code == "23505" {
								log.Printf("Entry exists: %s --- continuing...", queryB.IngredientID)
								continue
							}
						}
						log.Printf("Failed to create conversion at position: %d - %s", index, i.Name)
					}
				}
				bar.Add(1)
				continue
			}
		}

		queryA := database.CreateIngredientParams{
			Name: i.Name,
			ImageKey: ik,
		}

		ingredient, err := dbConn.CreateIngredient(ctx, queryA) 
		if err != nil {
			log.Printf("Failed on ingredient: %s - ERROR: %v", i.Name, err)
			log.Panic("Couldn't create ingredient during setup")
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

	log.Println("Database verification successful...")
	return nil
}
