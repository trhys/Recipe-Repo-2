package main

import (
	"time"

	"github.com/trhys/Recipe-Repo-2/internal/database"
)

type apiConfig struct {
	db		*database.Queries
	platform	string
	secret		string
	jwtDuration	time.Duration
}
