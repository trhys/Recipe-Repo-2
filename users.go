package main

import (
	"encoding/json"
	"fmt"
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
}

type userResponse struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string `json:"email"`
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
	}

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
	JWT		string `json:"token"`
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
		respondJSON(w, 201, loginResponse{
			ID: user.ID,
			Email: req.Email,
			JWT: token,
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

// Get User by ID
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
	}

	respondJSON(w, 200, res)
}
