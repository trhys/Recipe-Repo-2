package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/auth"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

func (cfg *apiConfig) handlerCreateShoppingList(w http.ResponseWriter, r *http.Request) {
	var req struct{
		Name	string `json:"name"`
		Recipes []struct{
			ID uuid.UUID `json:"id"`
			Quantity float32 `json:"quantity"`
		} `json:"recipes"`
	}


	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		respondFail(w, 400, "Failed to decode request body", err)
		return
	}

	// AUTH
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondFail(w, 400, "Invalid header", err)
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondFail(w, 401, "Invalid token", err)
		return
	}

	// Create list
	list, err := cfg.db.CreateShoppingList(r.Context(), database.CreateShoppingListParams{
		Name: req.Name,
		UserID: subject,
	})

	// Get ingredients
	ingredients := []database.AddToShoppingListParams{}
	for _, recipe := range req.Recipes {
		recipeIngredients, err := cfg.db.GetRecipesIngredients(r.Context(), recipe.ID)
		if err != nil {
			respondFail(w, 404, "Couldn't create list", err)
			return
		}

		for _, id := range recipeIngredients {
			ingredients = append(ingredients, database.AddToShoppingListParams{
				ShoppingListID: list.ID,
				IngredientID: id,
				Quantity: recipe.Quantity,
			})
		}	
	}

	// Add ings to list
	for _, ing := range ingredients {
		err := cfg.db.AddToShoppingList(r.Context(), ing)
		if err != nil {
			respondFail(w, 500, "something went wrong", err)
			return
		}
	}

	
	respondJSON(w, 204, nil)
}
