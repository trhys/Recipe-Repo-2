package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trhys/Recipe-Repo-2/internal/database"
	"github.com/trhys/Recipe-Repo-2/internal/auth"
)

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
		ID: uuid.New(),
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
