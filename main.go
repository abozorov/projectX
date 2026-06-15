package main

import (
	"context"
	"fmt"
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
	events "github.com/abozorov/projectX/internal/eventbus"
	"github.com/abozorov/projectX/internal/repo"
	requestQueue "github.com/abozorov/projectX/internal/request_queue"
	"github.com/abozorov/projectX/internal/service"
	"github.com/abozorov/projectX/pkg/db"
	"github.com/abozorov/projectX/pkg/logger"
	_ "github.com/lib/pq" // To register the driver.
	"go.uber.org/zap"
)

/*
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
*/
func main() {
	// get config
	cfg, err := config.NewConfig("internal/config/config.env")
	if err != nil {
		log.Fatal(err)
	}

	// for logger
	ctx, cancle := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// create logger
	logger, err := logger.NewLogger(true, cfg.AuditLogStorage)
	if err != nil {
		log.Println("Func main", zap.Error(err))
		return
	}

	// db connettion
	db, err := db.New(db.Options{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		logger.Error("main", zap.Error(err))
		return
	}

	// create event bus
	bus := events.NewBus(10)
	consumer.StartAuditConsumer(ctx, &wg, bus, logger)

	// register queue
	queue := requestQueue.NewQueueLimit(10)
	requestQueue.StartQueueConsumer(ctx, &wg, queue, logger)

	st := repo.NewpostgresRepo(db)
	service := service.NewUserService(st, bus)
	h := handlers.NewUserHandler(service, logger, queue)

	// create server
	router := handlers.NewRouter(h, queue)
	server := &http.Server{
		Addr:    cfg.HttpHost,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server started localhost:%s started", server.Addr))
		err = server.ListenAndServe()
		if err != nil {
			logger.Error("main", zap.Error(err))
			return
		}
	}()

	// gracefully shutdown
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
