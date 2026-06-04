package handlers

import (
	"net/http"

	requestQueue "github.com/abozorov/projectX/internal/request_queue"
)

type Router struct {
	*http.ServeMux
	QueueLimit *requestQueue.QueueLimit
}

func NewRouter(h *UserHandler, queue *requestQueue.QueueLimit) *Router {
	mux := http.NewServeMux()

	// user handlers
	mux.Handle("GET /users", http.HandlerFunc(h.GetUsers))
	mux.Handle("GET /user/{user_id}", http.HandlerFunc(h.GetUserByID))
	mux.Handle("POST /user", http.HandlerFunc(h.CreateUser))
	mux.Handle("PUT /user/{user_id}", http.HandlerFunc(h.UpdateUser))

	return &Router{
		mux,
		queue,
	}
}
