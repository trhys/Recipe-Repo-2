package viewmodel

import (
        "fmt"
        "time"

        "github.com/google/uuid"
        "github.com/trhys/Recipe-Repo-2/internal/database"
)

type Recipe struct {
    ID          uuid.UUID            `json:"id"`
    Title       string               `json:"title"`
    CreatedAt   time.Time            `json:"created_at"`
    UpdatedAt   *time.Time           `json:"updated_at,omitempty"`
    UserID      *uuid.UUID           `json:"user_id,omitempty"`
    Author      string               `json:"author,omitempty"`
    Description string               `json:"description,omitempty"`
    ImageURL    string               `json:"image_url,omitempty"`
    Ingredients	[]Ingredient	     `json:'ingredients,omitempty"`
}

type RecipeCardViewModel struct {
	Recipes	[]Recipe	`json:"recipes"`
}

func (builder *VMFactory) GenerateRecipeCardViewModel(recipes []database.Recipe) RecipeCardViewModel {
	model := RecipeCardViewModel{}
	for _, r := range recipes {
		model.Recipes = append(model.Recipes, Recipe {
			ID:		r.ID,
			Title:		r.Title,
			CreatedAt:	r.CreatedAt,
			UserID:		&r.UserID,
			Author:		r.Author,
			ImageURL:	fmt.Sprintf("%s/%s", builder.S3cdn, r.ImageKey),
		})
	}

	return model
}

func (builder *VMFactory) GenerateRecipeFullViewModel(r database.Recipe, i []Ingredient) Recipe {
	return Recipe {
		ID:		r.ID,
		Title:		r.Title,
		CreatedAt:	r.CreatedAt,
		UpdatedAt:	&r.UpdatedAt,
		UserID:		&r.UserID,
		Author:		r.Author,
		Description:	r.Description,
		ImageURL:	fmt.Sprintf("%s/%s", builder.S3cdn, r.ImageKey),
		Ingredients:	i,
	}
}
