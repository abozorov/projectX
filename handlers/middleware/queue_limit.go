package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	requestQueue "github.com/abozorov/projectX/internal/request_queue"
	"github.com/abozorov/projectX/package/errs"
)

func QueueLimit(queue *requestQueue.QueueLimit, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get client id
		cliID, err := strconv.Atoi(r.Header.Get("client_id"))
		if err != nil {
			log.Printf("[ERROR] Bad header: client_id is not an integer!")
			http.Error(w, errs.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("middlewware QueueLimit request cliID %d, len of queue %d\n", cliID, queue.QueueLen(cliID))

		// check quantity
		if queue.CheckQueue(cliID) {
			http.Error(w, errs.ErrTooManyRequests.Error(), http.StatusTooManyRequests)
			return
		}

		// add request to queue
		reqID := rand.Int()
		queue.AddClientRequest(cliID, reqID)
		r.Header.Add("request_id", strconv.Itoa(reqID))

		// next func
		next.ServeHTTP(w, r)
	})
}
