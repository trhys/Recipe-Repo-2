package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"html/template"
	"mime"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	"github.com/trhys/Recipe-Repo-2/internal/auth"
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
	Author		string `json:"author"`
	Description	string `json:"description"`
	ImageKey	string `json:"image_key"`
}

type appResponse struct {
	recipeResponse
	ImageURL	string `json:"image_url"`
}

func (cfg *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request) {
	// Request
	r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
	var req struct{
		Title string `json:"title"`
		UserID uuid.UUID `json:"user_id"`
		Ingredients []struct{
			Name		string `json:"name"`
			Quantity	float32 `json:"quantity"`
			Unit		string `json:"unit"`
		} `json:"ingredients"`
		Description string `json:"description"`
	}

	// Get request payload 
	jsonString := r.FormValue("payload")

	// Unmarshal JSON
	if err := json.Unmarshal([]byte(jsonString), &req); err != nil {
		respondFail(w, 500, "Failed to unmarshal payload", err)
		return
	}

	// Make sure user exists
	user, err := cfg.db.GetUser(r.Context(), req.UserID)
	if err != nil {
		respondFail(w, 404, "Invalid user id", err)
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

	// Request is valid - begin processing image file
	file, fileHeader, err := r.FormFile("image")
	key := uuid.New().String()
	if err == nil {
		defer file.Close()
	
		mediaType, _, err := mime.ParseMediaType(fileHeader.Header.Get("Content-Type"))
		if err != nil {
			respondFail(w, 401, "Couldn't parse media type", err)
			return
		}

		if mediaType != "image/jpeg" && mediaType != "image/png" {
			respondFail(w, 401, "Invalid media type", fmt.Errorf("Must be jpg or png. Got: %s", mediaType))
			return
		}

		tmp, err := os.CreateTemp("", "image_upload")
		if err != nil {
			respondFail(w, 500, "Something went wrong", err)
			return
		}
		defer os.Remove(tmp.Name())
		defer tmp.Close()

		_, fail := io.Copy(tmp, file)
		if fail != nil {
			respondFail(w, 500, "Couldn't save image", err)
			return
		}

		tmp.Seek(0, io.SeekStart)

		// Upload to s3
		if _, err := cfg.s3client.PutObject(r.Context(), &s3.PutObjectInput{
			Bucket: &cfg.s3bucket,
			Key: &key,
			Body: tmp,
			ContentType: &mediaType,
		}); err != nil {
			respondFail(w, 500, "Failed to upload to s3 bucket", err)
			return
		}

	} else if err != nil {
		if err != http.ErrMissingFile {
			respondFail(w, 500, "Something went wrong during upload", err)
			return
		}
	}

	// Query database
	query := database.CreateRecipeParams{
		Title: req.Title,
		UserID: user.ID,
		Description: req.Description,
		ImageKey: key,
	}

	recipe, err := cfg.db.CreateRecipe(r.Context(), query)
	if err != nil {
		respondFail(w, 404, "Couldn't create recipe", err)
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
	
	// Response
	res := recipeResponse{
		ID: recipe.ID,
		Title: recipe.Title,
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		UserID: recipe.UserID,
		Ingredients: ingredients,
		Author: user.Name,
		Description: recipe.Description,
		ImageKey: recipe.ImageKey,
	}

	respondJSON(w, 201, res)
}

// Get recipe by ID
// API endpoint
func (cfg *apiConfig) handlerGetRecipe(w http.ResponseWriter, r *http.Request) {
	requested := r.PathValue("recipe_id")
	recipe_id, err := uuid.Parse(requested)
	if err != nil {
		respondFail(w, 404, "Invalid recipe id", err)
		return
	}

	recipe, err := cfg.db.GetRecipe(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipe id", err)
		return
	}

	i, err := cfg.db.GetIngredientList(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find ingredients", err)
		return
	}

	author, err := cfg.db.GetName(r.Context(), recipe.UserID)
	if err != nil {
		respondFail(w, 404, "Couldn't find author", err)
		return
	}

	ingredients := []ingredient{}
	for _, ing := range i {
		ingredients = append(ingredients, ingredient{
			ID: ing.ID,
			Name: ing.Name,
			Quantity: ing.Quantity,
			Unit: ing.Unit,
			CreatedAt: ing.CreatedAt,
			UpdatedAt: ing.UpdatedAt,
			RecipeID: ing.RecipeID,
		})
	}

	res := recipeResponse{
		ID: recipe.ID,
		Title: recipe.Title,
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		UserID: recipe.UserID,
		Ingredients: ingredients,
		Author: author,
		Description: recipe.Description,
		ImageKey: recipe.ImageKey,
	}

	respondJSON(w, 200, res)
}

// App endpoint
func (cfg *apiConfig) appGetRecipe(w http.ResponseWriter, r *http.Request) {
	requested := r.PathValue("recipe_id")
	recipe_id, err := uuid.Parse(requested)
	if err != nil {
		respondFail(w, 404, "Invalid recipe id", err)
		return
	}

	recipe, err := cfg.db.GetRecipe(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipe id", err)
		return
	}

	i, err := cfg.db.GetIngredientList(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find ingredients", err)
		return
	}

	author, err := cfg.db.GetName(r.Context(), recipe.UserID)
	if err != nil {
		respondFail(w, 404, "Couldn't find author", err)
		return
	}

	ingredients := []ingredient{}
	for _, ing := range i {
		ingredients = append(ingredients, ingredient{
			ID: ing.ID,
			Name: ing.Name,
			Quantity: ing.Quantity,
			Unit: ing.Unit,
			CreatedAt: ing.CreatedAt,
			UpdatedAt: ing.UpdatedAt,
			RecipeID: ing.RecipeID,
		})
	}

	data := appResponse{
		recipeResponse: recipeResponse{
			ID: recipe.ID,
			Title: recipe.Title,
			CreatedAt: recipe.CreatedAt,
			UpdatedAt: recipe.UpdatedAt,
			UserID: recipe.UserID,
			Ingredients: ingredients,
			Author: author,
			Description: recipe.Description,
			ImageKey: recipe.ImageKey,
		},
		ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, recipe.ImageKey),
	}
	
	tmpl, _ := template.ParseFiles(filepath.Join("app", "templates", "recipe-viewer.html"))
	tmpl.Execute(w, data)
}	

// Get ten most recent recipes
type recipe struct{
	ID		uuid.UUID `json:"id"`
	Title		string `json:"title"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	UserID		uuid.UUID `json:"user_id"`
	Author		string `json:"author"`
	ImageKey	string `json:"image_key"`
	ImageURL	string `json:"image_url"`
}

type recipeList struct{
	Recipes []recipe `json:"recipes"`
}

func (cfg *apiConfig) handlerGetRecipeList(w http.ResponseWriter, r *http.Request) {
	recipes, err := cfg.db.GetRecipeList(r.Context())
	if err != nil {
		respondFail(w, 404, "Failed to retrieve recipe list", err)
		return
	}

	list := recipeList{}
	for _, rec := range recipes{
		author, err := cfg.db.GetName(r.Context(), rec.UserID)
		if err != nil {
			respondFail(w, 404, "Couldn't resolve author", err)
			return
		}

		list.Recipes = append(list.Recipes, recipe{
			ID: rec.ID,
			Title: rec.Title,
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			UserID: rec.UserID,
			Author: author,
			ImageKey: rec.ImageKey,
			ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, rec.ImageKey),
		})
	}

	respondJSON(w, 200, list)
}

// Get all recipes for user_id
type userRecipeList struct{
	Recipes []recipe `json:"recipes"`
	Name	string `json:"name"`
}

func (cfg *apiConfig) handlerGetUsersRecipes(w http.ResponseWriter, r *http.Request) {
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

	recipes, err := cfg.db.GetUsersRecipes(r.Context(), user.ID)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipes", err)
		return
	}

	list := userRecipeList{}
	for _, rec := range recipes {
		list.Recipes = append(list.Recipes, recipe{
			ID: rec.ID,
			Title: rec.Title,
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			UserID: rec.UserID,
			Author: user.Name,
			ImageKey: rec.ImageKey,
			ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, rec.ImageKey),
		})
	}
	list.Name = user.Name

	respondJSON(w, 200, list)
}

func (cfg *apiConfig) appGetUsersRecipes(w http.ResponseWriter, r *http.Request) {
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

	recipes, err := cfg.db.GetUsersRecipes(r.Context(), user.ID)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipes", err)
		return
	}

	list := userRecipeList{}
	for _, rec := range recipes {
		list.Recipes = append(list.Recipes, recipe{
			ID: rec.ID,
			Title: rec.Title,
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			Author: user.Name,
			ImageKey: rec.ImageKey,
			ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, rec.ImageKey),
		})
	}
	list.Name = user.Name

	tmpl, err := template.ParseFiles(filepath.Join("app", "templates", "user_page.html"))
	if err != nil {
		respondFail(w, 500, "Something went wrong", err)
		return
	}

	tmpl.Execute(w, list)
}
