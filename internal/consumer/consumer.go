package consumer

import (
	"context"
	"sync"

	events "github.com/abozorov/projectX/internal/eventbus"
	"github.com/abozorov/projectX/pkg/logger"
	"go.uber.org/zap"
)

func StartAuditConsumer(ctx context.Context, wg *sync.WaitGroup, bus *events.Bus, log *logger.Logger) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Info("Audit Consumer closed")
				return
			case event := <-bus.Subscribe():
				log.Audit.Info(
					event.Type,
					zap.Int("client_id", event.ClientId),
				)
			// default:
			// 	log.Info("") // log dropped event
			}
		}

	}()
}
