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

type ingredientResponse struct {
	ID		uuid.UUID `json:"id"`
	Name		string `json:"name"`
	Quantity 	float32 `json:"quantity"`
	Unit		string `json:"unit"`
}

// TODO: document ep usage for struct fields to clarify what is omitted in various endpoints

type recipe struct {
    ID          uuid.UUID            `json:"id"`
    Title       string               `json:"title"`
    CreatedAt   time.Time            `json:"created_at"`
    UpdatedAt   *time.Time           `json:"updated_at,omitempty"`
    UserID      *uuid.UUID           `json:"user_id,omitempty"`
    Author      string               `json:"author,omitempty"`
    Description string               `json:"description,omitempty"`
    ImageKey    string               `json:"image_key,omitempty"`
    Ingredients []ingredientResponse `json:"ingredients,omitempty"`
    ImageURL    string               `json:"image_url,omitempty"`
}

type recipeList struct{
	Recipes []recipe `json:"recipes"`
}
func (cfg *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request) {
	// Request
	r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
	var req struct{
		Title string `json:"title"`
		UserID uuid.UUID `json:"user_id"`
		Description string `json:"description"`
		Ingredients []struct{
			ID		uuid.UUID `json:"id"`
			Quantity 	float32 `json:"quantity"`
			Unit		string `json:"unit"`
		} `json:"ingredients"`

	}

	// Get request payload 
	jsonString := r.FormValue("payload")

	// Unmarshal JSON
	if err := json.Unmarshal([]byte(jsonString), &req); err != nil {
		respondFail(w, 500, "Failed to unmarshal payload", err)
		return
	}

	// Get username
	username, err := cfg.db.GetName(r.Context(), req.UserID)
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

	if subject != req.UserID {
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
		Author: username,
		UserID: req.UserID,
		Description: req.Description,
		ImageKey: key,
	}

	rec, err := cfg.db.CreateRecipe(r.Context(), query)
	if err != nil {
		respondFail(w, 404, "Couldn't create recipe", err)
		return
	}

	// Connect all ingredients
	for _, ing := range req.Ingredients {
		query := database.AddToRecipeParams{
			RecipeID: rec.ID,
			IngredientID: ing.ID,
			Quantity: ing.Quantity,
			Unit: ing.Unit,
		}

		_, err := cfg.db.AddToRecipe(r.Context(), query)
		if err != nil {
			respondFail(w, 500, "Failed to add ingredient", err)
			return
		}	
	}
	
	// Response
	res := recipe{
		ID: rec.ID,
		Title: rec.Title,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: &rec.UpdatedAt,
		UserID: &rec.UserID,
		Author: username,
		Description: rec.Description,
		ImageKey: rec.ImageKey,
	}

	respondJSON(w, 201, res)
}

// Get recipe by ID
func (cfg *apiConfig) handlerGetRecipe(w http.ResponseWriter, r *http.Request) {
	requested := r.PathValue("recipe_id")
	recipe_id, err := uuid.Parse(requested)
	if err != nil {
		respondFail(w, 404, "Invalid recipe id", err)
		return
	}

	rec, err := cfg.db.GetRecipe(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find recipe id", err)
		return
	}

	i, err := cfg.db.GetIngredientList(r.Context(), recipe_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find ingredients", err)
		return
	}

	author, err := cfg.db.GetName(r.Context(), rec.UserID)
	if err != nil {
		respondFail(w, 404, "Couldn't find author", err)
		return
	}

	ingredients := []ingredientResponse{}
	for _, ing := range i {
		ingredients = append(ingredients, ingredientResponse{
			ID: ing.IngredientID,
			Name: ing.Name,
			Quantity: ing.Quantity,
			Unit: ing.Unit,
		})
	}

	res := recipe{
		ID: rec.ID,
		Title: rec.Title,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: &rec.UpdatedAt,
		UserID: &rec.UserID,
		Ingredients: ingredients,
		Author: author,
		Description: rec.Description,
		ImageKey: rec.ImageKey,
		ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, rec.ImageKey),
	}

	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, 200, res)
		return
	}

	tmpl, _ := template.ParseFiles(filepath.Join("app", "templates", "recipe-viewer.html"))
	tmpl.Execute(w, res)
}

// Get ten most recent recipes
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
			UpdatedAt: &rec.UpdatedAt,
			UserID: &rec.UserID,
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
			UpdatedAt: &rec.UpdatedAt,
			UserID: &rec.UserID,
			Author: user.Name,
			ImageKey: rec.ImageKey,
			ImageURL: fmt.Sprintf("%s/%s", cfg.s3cdn, rec.ImageKey),
		})
	}
	list.Name = user.Name

	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, 200, list)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("app", "templates", "user_page.html"))
	if err != nil {
		respondFail(w, 500, "Something went wrong", err)
		return
	}

	tmpl.Execute(w, list)
}
