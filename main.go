package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/abozorov/projectX/handlers"
	"github.com/abozorov/projectX/middleware"
	"github.com/abozorov/projectX/models"
	"github.com/abozorov/projectX/storage"
)

func initFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	users := []models.User{}
	err = json.NewDecoder(file).Decode(&users)

	if err != nil {
		file.WriteString("[]")
	}
	return nil
}

func main() {

	st := &storage.UserStorage{
		FileName: "data/users.json",
	}

	err := initFile(st.FileName)
	if err != nil {
		fmt.Println(err)
		return
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
