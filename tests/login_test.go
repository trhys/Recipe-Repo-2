package main

import (
	"testing"
	"bytes"
	"net/http"
)

type tests struct {
	cases [][]byte
}

func TestLogin(t *testing.T) {
	http.Post("http://localhost:8080/api/reset", "application/json", nil)
	createUserReq := tests{
		cases: [][]byte{
			[]byte(`{"email": "test1", "password": "test1password"}`),
			[]byte(`{"email": "test2", "password": "theothertest123"}`),
		},
	}

	for _, req := range createUserReq.cases {
		_, err := http.Post("http://localhost:8080/api/new_user", "application/json", bytes.NewBuffer(req))
		if err != nil {
			t.Fatalf("error on request: %v", err)
		}

		res, err := http.Post("http://localhost:8080/api/login", "application/json", bytes.NewBuffer(req))
		if err != nil {
			t.Fatalf("error on request: %v", err)
		}

		if res.StatusCode == 401 {
			t.Fatalf("failed to authenticate")
		} else if res.StatusCode != 201 {
			t.Fatalf("something went wrong in authentication")
		}
	}
}	
