package handlers

import (
	"net/http"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(h *UserHandler) *Router {
	mux := http.NewServeMux()

	// user handlers
	mux.Handle("GET /users", http.HandlerFunc(h.GetUsers))
	mux.Handle("GET /user/{user_id}", http.HandlerFunc(h.GetUserByID))
	mux.Handle("POST /user", http.HandlerFunc(h.CreateUser))
	mux.Handle("PUT /user/{user_id}", http.HandlerFunc(h.UpdateUser))

	return &Router{
		mux,
	}
}
