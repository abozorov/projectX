package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/abozorov/projectX/handlers"
	"github.com/abozorov/projectX/handlers/middleware"
	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/internal/service"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/package/logger"
)

var (
	dataFile = "data/users.json"
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

	logger := logger.NewLogger(true)
	st := storage.NewUserStorage(dataFile)
	err := initFile(dataFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	service := service.NewUserService(st)
	h := handlers.NewUserHandler(service, logger)
	mux := handlers.NewRouter(h)

	handler := middleware.Logging(middleware.Auth(mux))
	fmt.Println("Server localhost:8080 started")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
		return
	}
}
