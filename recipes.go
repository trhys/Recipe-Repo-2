package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

type ingredient struct {
	ID		uuid.UUID `json:"id"`
	Name		string `json:"name"`
	Quantity	float32 `json:"quantity"`
	Unit		string `json:"unit"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	RecipeID	uuid.UUID `json:"recipe_id"`
}

type recipeResponse struct {
	ID		uuid.UUID `json:"id"`
	Title		string `json:"title"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	UserID		uuid.UUID `json:"user_id"`
	Ingredients	[]ingredient `json:"ingredients"`
}

func (cfg *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct{
		Title string `json:"title"`
		UserID uuid.UUID `json:"user_id"`
		Ingredients []struct{
			Name		string `json:"name"`
			Quantity	float32 `json:"quantity"`
			Unit		string `json:"unit"`
		} `json:"ingredients"`
	}
	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 500, "Failed to decode request body", err)
		return
	}

	query := database.CreateRecipeParams{
		Title: req.Title,
		UserID: req.UserID,
	}

	recipe, err := cfg.db.CreateRecipe(r.Context(), query)
	if err != nil {
		respondFail(w, 404, "Couldn't create recipe: does user_id exist?", err)
		return
	}

	ingredients := []ingredient{}
	for _, ing := range req.Ingredients {
		query := database.CreateIngredientParams{
			Name: ing.Name,
			Quantity: ing.Quantity,
			Unit: ing.Unit,
			RecipeID: recipe.ID,
		}

		i, err := cfg.db.CreateIngredient(r.Context(), query)
		if err != nil {
			respondFail(w, 500, "Failed to create ingredient", err)
			return
		}

		ingredients = append(ingredients, ingredient{
			ID: i.ID,
			Name: i.Name,
			Quantity: i.Quantity,
			Unit: i.Unit,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			RecipeID: i.RecipeID,
		})
	}
	
	res := recipeResponse{
		ID: recipe.ID,
		Title: recipe.Title,
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		UserID: recipe.UserID,
		Ingredients: ingredients,
	}

	respondJSON(w, 201, res)
}
