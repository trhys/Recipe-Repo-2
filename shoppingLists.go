package main

import (
	"encoding/json"
	"html/template"
	"path/filepath"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/auth"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

type shoppingList struct{
	ID		uuid.UUID 	`json:"id"`
	Name		string		`json:"name"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
} 

func (cfg *apiConfig) handlerCreateShoppingList(w http.ResponseWriter, r *http.Request) {
	var req struct{
		Name	string `json:"name"`
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

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 400, "Failed to decode request body", err)
		return
	}

	// Create list
	list, err := cfg.db.CreateShoppingList(r.Context(), database.CreateShoppingListParams{
		Name: req.Name,
		UserID: subject,
	})

	respondJSON(w, 200, list)
}

func (cfg *apiConfig) handlerAddToShoppingList(w http.ResponseWriter, r *http.Request) {
	var req struct{
		ShoppingListID	uuid.UUID `json:"shopping_list_id"`
		RecipeID	uuid.UUID `json:"recipe_id"`
		Quantity	int32	  `json:"quantity"`
	}

	// AUTH
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondFail(w, 400, "Invalid header", err)
                return
        }

        if _, err := auth.ValidateJWT(token, cfg.secret); err != nil {
                respondFail(w, 401, "Invalid token", err)
                return
        }

	// Decode body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 400, "Failed to decode request body", err)
		return
	}

	// Get ingredients from recipe id
	ingredientList, err := cfg.db.GetIngredientList(r.Context(), req.RecipeID)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipe", err)
		return
	}

	for _, ingredient := range ingredientList {
		err := cfg.db.AddToShoppingList(r.Context(), database.AddToShoppingListParams{
			ShoppingListID: req.ShoppingListID,
			IngredientID: ingredient.IngredientID,
			Quantity: ingredient.Quantity,
		})
		if err != nil {
			respondFail(w, 500, "Something went wrong", err)
			return
		}
	}

	// Link recipe to list
	if err := cfg.db.AddRecipeToList(r.Context(), database.AddRecipeToListParams{
		ShoppingListID: req.ShoppingListID,
		RecipeID: req.RecipeID,
		Quantity: req.Quantity,
	}); err != nil {
		respondFail(w, 500, "something went wrong", err)
		return
	}

	respondJSON(w, 204, nil)
}


func (cfg *apiConfig) handlerGetShoppingList(w http.ResponseWriter, r *http.Request) {
	val := r.PathValue("shopping_list_id")
	listID, err := uuid.Parse(val) 
	if err != nil {
		respondFail(w, 404, "invalid uuid", err)
		return
	}

	shoppingList, err := cfg.db.GetShoppingList(r.Context(), listID)
	if err != nil {
		respondFail(w, 404, "Couldn't find shopping list", err)
		return
	}

	type res struct{
		ID	 uuid.UUID		`json:"id"`
		Name	 string 		`json:"name"`
		Created  time.Time		`json:"created_at"`
		Recipes  []recipe 		`json:"recipes"`
		Quantity map[uuid.UUID]int32	`json:"quantity"`
	}

	shoppingListRecipes, err := cfg.db.GetRecipesFromList(r.Context(), shoppingList.ID)
	if err != nil {
		respondFail(w, 404, "couldnt find recipes from list", err)
		return
	}

	response := res{
		ID: shoppingList.ID,
		Name: shoppingList.Name,
		Created: shoppingList.CreatedAt,
		Quantity: make(map[uuid.UUID]int32),
	}

	for _, r := range shoppingListRecipes {
		response.Recipes = append(response.Recipes, recipe{
			ID: r.ID,
			Title: r.Title,
			Author: r.Author,
			UserID: &r.UserID,
		})

		response.Quantity[r.ID] = r.Quantity
	}

	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, 200, response)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("app", "templates", "shopping_list.html"))
        if err != nil {
                respondFail(w, 500, "Something went wrong", err)
                return
        }

        tmpl.Execute(w, response)
}

func (cfg *apiConfig) handlerGetUsersShoppingLists(w http.ResponseWriter, r *http.Request) {
	val := r.PathValue("user_id")
	id, err := uuid.Parse(val)
        if err != nil {
                respondFail(w, 404, "Invalid uuid", err)
                return
        }

        user, err := cfg.db.GetUser(r.Context(), id)
        if err != nil {
                respondFail(w, 404, "Couldn't find user", err)
                return
        }

	// Authorization
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondFail(w, 401, "Failed to retrieve bearer token", err)
                return
        }

        subject, err := auth.ValidateJWT(token, cfg.secret)
        if err != nil {
                respondFail(w, 401, "Couldn't validate token", err)
                return
        }

        if subject != user.ID {
                respondFail(w, 401, "Unauthorized access", nil)
                return
        }

	// Authorized - Get lists

	lists, err := cfg.db.GetUserLists(r.Context(), user.ID)
	if err != nil {
		respondFail(w, 404, "Couldn't retrieve user's shopping lists", err)
		return
	}

	var res struct{
		Name	string `json:"name"`
		Lists	[]shoppingList `json:"shopping_lists"`
	}

	res.Name = user.Name

	for _, list := range lists {
		res.Lists = append(res.Lists, shoppingList{
			ID: list.ID,
			Name: list.Name,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		})
	}

	if r.Header.Get("Accept") == "application/json" {
                respondJSON(w, 200, res)
                return
        }

        tmpl, err := template.ParseFiles(filepath.Join("app", "templates", "user_shopping_lists.html"))
        if err != nil {
                respondFail(w, 500, "Something went wrong", err)
                return
        }

        tmpl.Execute(w, res)
}

func (cfg *apiConfig) handlerPrintList(w http.ResponseWriter, r *http.Request) {
	// AUTH
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondFail(w, 400, "Invalid header", err)
                return
        }

        if userID, err := auth.ValidateJWT(token, cfg.secret); err != nil {
                respondFail(w, 401, "Invalid token", err)
                return
        }

	var req struct{
		ID uuid.UUID `json:"id"`
	}

	// Decode body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 400, "Failed to decode request body", err)
		return
	}



