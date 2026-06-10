package handlers

import (
	"net/http"

	"github.com/abozorov/projectX/handlers/middleware"
	requestQueue "github.com/abozorov/projectX/internal/request_queue"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(h *UserHandler, queue *requestQueue.QueueLimit) *Router {
	userMux := http.NewServeMux()

	// user handlers
	userMux.Handle("GET /users", http.HandlerFunc(h.GetUsers))
	userMux.Handle("GET /user/{user_id}", http.HandlerFunc(h.GetUserByID))
	userMux.Handle("POST /user", http.HandlerFunc(h.CreateUser))
	userMux.Handle("PUT /user", http.HandlerFunc(h.UpdateUser))
	userMux.Handle("DELETE /user/{user_id}", http.HandlerFunc(h.DeleteUser))
	userMux.Handle("/", middleware.QueueLimit(queue, middleware.Logging(middleware.Auth(userMux))))

	// login
	userMux.Handle("POST /login", middleware.QueueLimit(queue, middleware.Logging(http.HandlerFunc(h.Login))))

	return &Router{
		userMux,
	}
}
