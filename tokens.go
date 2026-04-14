package main

import (
	"net/http"

	"github.com/trhys/Recipe-Repo-2/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondFail(w, 401, "Invalid header", err)
		return
	}

	user, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondFail(w, 401, "Invalid token", err)
		return
	}
	
	jwt, err := auth.MakeJWT(user, cfg.secret, cfg.jwtDuration)
	if err != nil {
		respondFail(w, 500, "Failed to write JWT", err)
		return
	}

	type res struct{
		Token string `json:"token"`
	}

	resp := res{
		Token: jwt,
	}

	respondJSON(w, 200, resp)
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondFail(w, 401, "Invalid header", err)
		return
	}

	if err := cfg.db.RevokeToken(r.Context(), token); err != nil {
		respondFail(w, 401, "Invalid token", err)
		return
	}

	respondJSON(w, 204, nil)
}
