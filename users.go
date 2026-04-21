package main

import (
	"encoding/json"
	"path/filepath"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	"github.com/trhys/Recipe-Repo-2/internal/auth"
)

// Create user
type createUserRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
	Name		string `json:"name"`
}

type userResponse struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string `json:"email"`
	Name		string `json:"name"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req createUserRequest

	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 500, "Failed to decode response body", err)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondFail(w, 500, "Failed to hash password", err)
		return
	}

	query := database.CreateUserParams{
		Email: req.Email,
		HashedPw: hash,
		Name: req.Name,
	}

	user, err := cfg.db.CreateUser(r.Context(), query)
	if err != nil {
		respondFail(w, 500, "Database error: failed to create user", err)
		return
	}

	res := userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
		Name: user.Name,
	}

	log.Printf("User created with email: %s", res.Email)
	respondJSON(w, 201, res)
}


// User login
type loginRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

type loginResponse struct {
	ID		uuid.UUID `json:"id"`
	Email		string `json:"email"`
	Username	string `json:"username"`
	JWT		string `json:"token"`
	RT		string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req loginRequest

	if err := decoder.Decode(&req); err != nil {
		respondFail(w, 500, "Failed to decode response body", err)
		return
	}

	user, err := cfg.db.GetUserHash(r.Context(), req.Email)
	if err != nil {
		respondFail(w, 404, "User not found", err)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, user.HashedPw)
	if err != nil {
		respondFail(w, 500, "Something went wrong during authentication", err)
		return
	}

	if match {
		token, err := auth.MakeJWT(user.ID, cfg.secret, cfg.jwtDuration)
		if err != nil{
			respondFail(w, 500, "Failed to generate token", err)
			return
		}

		refresh_token := auth.MakeRefreshToken()
		if _, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			ID: refresh_token,
			UserID: user.ID,
		}); err != nil {
			respondFail(w, 500, "Couldn't generate refresh token", err)
			return
		}

		respondJSON(w, 201, loginResponse{
			ID: user.ID,
			Email: req.Email,
			Username: user.Name,
			JWT: token,
			RT: refresh_token,
		})
		return
	} else {
		respondFail(w, 401, "Invalid username or password", fmt.Errorf("Invalid username or password"))
		return
	}
}

// Reset users in dev mode
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondFail(w, 401, "Must be in dev mode to reset users!", nil)
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondFail(w, 500, "Something went wrong during users reset", err)
		return
	}

	log.Print("Reset database...")

	respondJSON(w, 201, nil)
}

// Get User by ID without recipe or shopping list info
func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")
	user_id, err := uuid.Parse(id)
	if err != nil {
		respondFail(w, 404, "Invalid user id", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), user_id)
	if err != nil {
		respondFail(w, 404, "Couldn't find user id", err)
		return
	}

	res := userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Name: user.Name,
	}

	respondJSON(w, 200, res)
}

// Get user profile - public
type userRecipeList struct{
        Recipes []recipe `json:"recipes"`
        Name    string `json:"name"`
}

func (cfg *apiConfig) handlerGetUserProfile(w http.ResponseWriter, r *http.Request) {
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
