package requestQueue

import (
	"context"
	"sync"

	"github.com/abozorov/projectX/pkg/logger"
)

func StartQueueConsumer(ctx context.Context, wg *sync.WaitGroup, queue *QueueLimit, log *logger.Logger) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Info("Queue Consumer closed")
				return
			case req := <-queue.Subscribe():
				queue.DeleteRequest(req.clientId, req.requestId)
			default:
				log.Info("")
			}
		}
	}()
}
