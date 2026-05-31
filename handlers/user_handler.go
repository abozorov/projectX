package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abozorov/projectX/models"
	"github.com/abozorov/projectX/package/errs"
	"github.com/abozorov/projectX/package/logger"
	"github.com/abozorov/projectX/storage"
	"go.uber.org/zap"
)

type UserHandler struct {
	storage *storage.UserStorage
}

func NewUserHandler(storage *storage.UserStorage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

var (
	log = logger.NewLogger(false)
)

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Info("Start func GetUsers")

	// load all
	users, err := h.storage.GetAll()
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Error("Func GetUsers", zap.String("error", err.Error()))
			http.Error(w, errs.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Error("Func GetUsers", zap.String("error", err.Error()))
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	// write request
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Error("Func GetUsers", zap.String("error", err.Error()))
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	log.Info("Start func GetUserByID", zap.String("user_id", r.PathValue("user_id")))

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		log.Error("Func GetUserByID", zap.String("error", err.Error()))
		http.Error(w, errs.ErrInvalidUserId.Error(), http.StatusBadRequest)
		return
	}

	// check id
	if id < 0 {
		log.Error("Func GetUserByID", zap.String("error", errs.ErrBadRequestQuery.Error()))
		http.Error(w, errs.ErrBadRequestQuery.Error(), http.StatusBadRequest)
		return
	}

	// get by id
	user, err := h.storage.GetByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrUserIDNotFound) {
			log.Error("Func GetUserByID", zap.String("error", err.Error()))
			http.Error(w, errs.ErrUserIDNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Error("Func GetUserByID", zap.String("error", err.Error()))
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
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
	log.Info("Start func CreateUser")

	// get user
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error("Func CreateUser", zap.String("error", err.Error()))
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}

	// id validation
	if user.ID < 0 {
		log.Error("Func CreateUser", zap.String("error", errs.ErrBadRequestBody.Error()))
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}

	// creating
	err = h.storage.Create(user)
	if err != nil {
		log.Error("Func CreateUser", zap.String("error", err.Error()))
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("User Created"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Start func UpdateUser", zap.String("user_id", r.PathValue("user_id")))

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		log.Error("Func UpdateUser", zap.String("error", err.Error()))
		http.Error(w, errs.ErrInvalidUserId.Error(), http.StatusBadRequest)
		return
	}

	// get user
	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error("Func UpdateUser", zap.String("error", err.Error()))
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusBadRequest)
		return
	}

	// check path id
	if id < 0 || user.ID < 0 || id != user.ID {
		log.Error("Func UpdateUser", zap.String("error", errs.ErrBadRequest.Error()))
		http.Error(w, errs.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	// updating
	err = h.storage.Update(id, user)
	if err != nil {
		log.Error("Func UpdateUser", zap.String("error", err.Error()))
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("User Updated"))
}
