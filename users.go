package main

import (
	"encoding/json"
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

type createUserResponse struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
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

	res := createUserResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
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
		respondJSON(w, 201, loginResponse{
			ID: user.ID,
			Email: req.Email,
		})
		return
	} else {
		respondJSON(w, 401, nil)
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

	respondJSON(w, 201, nil)
}
