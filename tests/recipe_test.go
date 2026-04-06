package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestRecipe(t *testing.T) {
	// make user 
	res, _ := http.Post("http://localhost:8080/api/new_user", "application/json",
		bytes.NewBuffer([]byte(`{"email": "recipe_tester", "password": "aw3s0m3r3c1p3"}`)))
	
	decoder := json.NewDecoder(res.Body)
	var user struct {
		ID uuid.UUID `json:"id"`
	}
	if err := decoder.Decode(&user); err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// make recipes
	cases := [][]byte{
		[]byte(fmt.Sprintf(`{"title": "an-awesome-recipe", "user_id": "%s"}`, user.ID)),
	}

	for _, req := range cases {
		res, err := http.Post("http://localhost:8080/api/new_recipe", "application/json", bytes.NewBuffer(req))
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		decoder := json.NewDecoder(res.Body)
		var recipe struct {
			ID	uuid.UUID `json:"id"`
			Title	string `json:"title"`
			UserID	uuid.UUID `json:"user_id"`
		}
		if err := decoder.Decode(&recipe); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if recipe.UserID != user.ID {
			t.Fatalf("Expected recipe with UserID: %s --- Got %s", user.ID, recipe.UserID)
		}
	}
}
