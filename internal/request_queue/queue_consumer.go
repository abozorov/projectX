package requestQueue

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
)

func StartQueueConsumer(ctx context.Context, wg *sync.WaitGroup, queue *QueueLimit) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Println("Queue Consumer closed")
				return
			case req := <-queue.Subscribe():
				idS := strings.Split(req, ":")
				cliID, _ := strconv.Atoi(idS[0])
				reqID, _ := strconv.Atoi(idS[1])

				queue.DeleteRequest(cliID, reqID)
				// default:
				// 	// log dropped event
				// 	log.Println("dropped event")
			}
		}
	}()
}
