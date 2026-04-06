package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

type recipeResponse struct {
	ID		uuid.UUID `json:"id"`
	Title		string `json:"title"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	UserID		uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct{
		Title string `json:"title"`
		UserID uuid.UUID `json:"user_id"`
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
	
	res := recipeResponse{
		ID: recipe.ID,
		Title: recipe.Title,
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		UserID: recipe.UserID,
	}

	respondJSON(w, 201, res)
}
