package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abozorov/projectX/internal/models"
	"github.com/abozorov/projectX/internal/service"
	"github.com/abozorov/projectX/package/errs"
	"github.com/abozorov/projectX/package/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
	log     *logger.Logger
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUserHandler(service *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     logger,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func GetUsers")

	// load all
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		errDistributor(err, w)
		h.log.Error("Func GetUsers", zap.String("error", err.Error()))
		return
	}

	// transform models.User -> user
	resp := make([]user, 0, len(users))
	for _, v := range users {
		resp = append(resp, user{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	// write request
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("Func GetUsers", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func GetUserByID", zap.String("user_id", r.PathValue("user_id")))

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		h.log.Error("Func GetUserByID", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequest, w)
		return
	}

	// get by id
	usr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("Func GetUserByID", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// transform models.User -> user
	resp := user{
		ID:   usr.ID,
		Name: usr.Name,
	}

	// write response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("Func GetUserByID", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func CreateUser")

	// get user
	usr := user{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.log.Error("Func CreateUser", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequestBody, w)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Create(r.Context(), models.User{
		ID:   usr.ID,
		Name: usr.Name,
	})
	if err != nil {
		h.log.Error("Func CreateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
	w.Write([]byte("User Created"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func UpdateUser", zap.String("user_id", r.PathValue("user_id")))

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequestQuery, w)
		return
	}

	// get user
	usr := user{}
	err = json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// updating
	err = h.service.Update(r.Context(), id, models.User{
		ID:   usr.ID,
		Name: usr.Name,
	})
	if err != nil {
		h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
	w.Write([]byte("User Updated"))
}
