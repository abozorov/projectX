package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abozorov/projectX/errs"
	"github.com/abozorov/projectX/models"
	"github.com/abozorov/projectX/storage"
)

type UserHandler struct {
	Storage *storage.UserStorage
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// load all
	users, err := h.Storage.GetAll()
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write request
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// check id
	if id < 0 {
		http.Error(w, errs.ErrBadRequestQuery.Error(), http.StatusBadRequest)
		return
	}

	// get by id
	user, err := h.Storage.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrUserIDNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write request
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// get user
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// id validation
	if user.ID < 0 {
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}

	// creating
	err = h.Storage.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("User Created"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// get user
	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check path id
	if id < 0 || user.ID < 0 || id != user.ID {
		http.Error(w, errs.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}
	
	// updating
	err = h.Storage.Update(id, user)
	if err != nil {
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("User Updated"))
}
