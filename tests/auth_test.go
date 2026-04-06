package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func AuthTest(t *testing.T) {
	res, err := http.Post("http://localhost:8080/api/new_user", "application/json", bytes.NewBuffer([]byte(`{"email": "authorizeduser@test.com", "password": "password1"}`)))
	decoder := json.NewDecoder(res.Body)
	var user struct{
		ID	uuid.UUID `json:"id"`
		Email	string `json:"email"`
	}

	if err := decoder.Decode(&user); err != nil {
		t.Fatal("Failed to decode new user response")
	}

	loginRes, err := http.Post("http://localhost:8080/api/login", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"email": "%s", "password": "password1"}`, user.Email))))
	loginDecoder := json.NewDecoder(loginRes.Body)
	var loggedIn struct{
		ID	uuid.UUID `json:"id"`
		Email	string `json:"email"`
		Token	string `json:"token"`
	}

	if err := loginDecoder.Decode(&loggedIn); err != nil {
		t.Fatal("Failed to decode login response")
	}

	if loggedIn.ID != user.ID || loggedIn.Email != user.Email {
		t.Fatal("Expected same user information between requests!")
	}

	authorizedRequest := []byte(fmt.Sprintf(`{"title": "a recipe", "user_id": "%s", "ingredients": [{"name": "apple", "quantity": 2, "unit": "whole"}, {"name": "cinnamon", "quantity": 3.5, "unit": "teaspoons"}]}`, loggedIn.ID))

	// Test good request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/new_recipe", bytes.NewBuffer(authorizedRequest))
	if err != nil {
		t.Fatalf("Something went wrong: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loggedIn.Token))

	client := &http.Client{}
	authResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if authResp.StatusCode == 401 {
		t.Fatal("Bad status code: expected to be authorized")
	}

	if authResp.StatusCode != 201 {
		t.Fatalf("Got non 401 bad status code: %d", authResp.StatusCode)
	}

	// Test request without token
	badReq, err := http.NewRequest("POST", "http://localhost:8080/api/new_recipe", bytes.NewBuffer(authorizedRequest))
	if err != nil {
		t.Fatalf("Something went wrong: %v", err)
	}

	badReq.Header.Set("Content-Type", "application/json")

	unauthResp, err := client.Do(badReq)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if unauthResp.StatusCode != 401 {
		t.Fatal("Bad status code: expected to be UNauthorized")
	}

}
