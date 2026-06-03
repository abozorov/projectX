package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/abozorov/projectX/handlers"
	"github.com/abozorov/projectX/handlers/middleware"
	"github.com/abozorov/projectX/internal/consumer"
	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/internal/service"
	events "github.com/abozorov/projectX/internal/service/eventbus"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/package/logger"
	"go.uber.org/zap"
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

	ctx, cancle := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}

	logger, err := logger.NewLogger(true)
	if err != nil {
		logger.Error("Func main", zap.Error(err))
		return
	}

	bus := events.NewBus(10)

	consumer.StartAuditConsumer(ctx, &wg, bus, logger)

	st := storage.NewUserStorage(dataFile)

	err = initFile(dataFile)
	if err != nil {
		logger.Error("Func main", zap.Error(err))
		return
	}

	service := service.NewUserService(st, bus)

	h := handlers.NewUserHandler(service, logger)

	handler := middleware.Logging(middleware.Auth(handlers.NewRouter(h)))
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		log.Printf("Server started localhost:%s started", server.Addr)
		err = server.ListenAndServe()
		if err != nil {
			logger.Error("Func main", zap.Error(err))
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	cancle()

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	wg.Wait()
	logger.Info("Server shutdown completed")
}
