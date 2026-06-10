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
	"github.com/abozorov/projectX/internal/config"
	"github.com/abozorov/projectX/internal/consumer"
	"github.com/abozorov/projectX/internal/models"
	requestQueue "github.com/abozorov/projectX/internal/request_queue"
	"github.com/abozorov/projectX/internal/service"
	events "github.com/abozorov/projectX/internal/service/eventbus"
	"github.com/abozorov/projectX/internal/storage"
	"github.com/abozorov/projectX/pkg/logger"
	"go.uber.org/zap"
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
	//
	cnf, err := config.NewConfig("internal/config/config.env")
	if err != nil {
		log.Fatal(err)
	}

	// for logger
	loggerCtx, loggerCancle := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// create logger
	logger, err := logger.NewLogger(true, cnf.AuditLogStorage)
	if err != nil {
		logger.Error("Func main", zap.Error(err))
		return
	}

	// create event bus
	bus := events.NewBus(10)
	consumer.StartAuditConsumer(loggerCtx, &wg, bus, logger)

	err = initFile(cnf.Storage)
	if err != nil {
		logger.Error("Func main", zap.Error(err))
		return
	}

	// start queue consumer
	// queueCtx, queueCancle := context.WithCancel(context.Background())

	queue := requestQueue.NewQueueLimit(10)
	requestQueue.StartQueueConsumer(loggerCtx, &wg, queue)

	st := storage.NewUserStorage(cnf.Storage)
	service := service.NewUserService(st, bus)
	h := handlers.NewUserHandler(service, logger, queue)

	// create server
	router := handlers.NewRouter(h, queue)
	server := &http.Server{
		Addr:    cnf.HttpHost,
		Handler: router,
	}

	go func() {
		log.Printf("Server started localhost:%s started", server.Addr)
		err = server.ListenAndServe()
		if err != nil {
			logger.Error("Func main", zap.Error(err))
			return
		}
	}()

	// gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	loggerCancle()
	// queueCancle()

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	wg.Wait()
	logger.Info("Server shutdown completed")
}
