package main

import (
	"fmt"
	"net/http"

	"github.com/abozorov/projectX/handlers"
	"github.com/abozorov/projectX/middleware"
	"github.com/abozorov/projectX/storage"
)

func main() {

	st := &storage.UserStorage{
		FileName: "data/users.json",
	}

	h := &handlers.UserHandler{
		Storage: st,
	}

	mux := http.NewServeMux()

	// GET /users
	mux.Handle("GET /users", http.HandlerFunc(h.GetUsers))

	// GET /users/{id}
	mux.Handle("GET /users/{user_id}", http.HandlerFunc(h.GetUserByID))

	// POST /users
	mux.Handle("POST /users", http.HandlerFunc(h.CreateUser))

	// PUT /users/{id}
	mux.Handle("PUT /users/{user_id}", http.HandlerFunc(h.UpdateUser))

	handler := middleware.Logging(middleware.Auth(mux))
	fmt.Println("Server localhost:8080 started")
	http.ListenAndServe(":8080", handler)
}
