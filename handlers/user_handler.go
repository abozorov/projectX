package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abozorov/projectX/internal/models"
	requestQueue "github.com/abozorov/projectX/internal/request_queue"
	"github.com/abozorov/projectX/internal/service"
	"github.com/abozorov/projectX/pkg/errs"
	"github.com/abozorov/projectX/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
	log     *logger.Logger
	queue   *requestQueue.QueueLimit
}

type updateUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type getUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUserHandler(service *service.UserService, logger *logger.Logger, queue *requestQueue.QueueLimit) *UserHandler {
	return &UserHandler{
		service: service,
		log:     logger,
		queue:   queue,
	}
}

func (h *UserHandler) doneRequest(r *http.Request) {
	// request done
	cliID, _ := strconv.Atoi(r.Header.Get("client_id"))
	reqID, _ := strconv.Atoi(r.Header.Get("request_id"))

	h.queue.Publish(cliID, reqID)
}

func (h *UserHandler) updateContext(ctx context.Context, r *http.Request) context.Context {
	cliID, _ := strconv.Atoi(r.Header.Get("client_id"))
	reqID, _ := strconv.Atoi(r.Header.Get("request_id"))

	ctx, _ = context.WithDeadline(ctx, h.queue.GetEndTime(cliID, reqID))
	return ctx
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// defer request done
	defer h.doneRequest(r)

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		h.log.Error("Func DeleteUser", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequest, w)
		return
	}

	// delete user
	err = h.service.DeleteUser(h.updateContext(r.Context(), r), id)
	if err != nil {
		h.log.Error("Func DeleteUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
	w.Write([]byte("user deleted"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// request done
	defer h.doneRequest(r)

	// get user
	usr := updateUser{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.log.Error("Func Login", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// check user
	autorizationKey, err := h.service.Login(h.updateContext(r.Context(), r), models.User{
		ID:       usr.ID,
		Password: usr.Password,
	})
	if err != nil {
		h.log.Error("Func Login", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// return autorization code
	w.Write([]byte("key: " + autorizationKey))
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func GetUsers")

	// defer request done
	defer h.doneRequest(r)

	// load all
	users, err := h.service.GetAll(h.updateContext(r.Context(), r))
	if err != nil {
		errDistributor(err, w)
		h.log.Error("Func GetUsers", zap.String("error", err.Error()))
		return
	}

	// transform models.User -> user
	resp := make([]getUser, 0, len(users))
	for _, v := range users {
		resp = append(resp, getUser{
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

	// defer request done
	defer h.doneRequest(r)

	// check path
	id, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		h.log.Error("Func GetUserByID", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequest, w)
		return
	}

	// get by id
	usr, err := h.service.GetByID(h.updateContext(r.Context(), r), id)
	if err != nil {
		h.log.Error("Func GetUserByID", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// transform models.User -> user
	resp := getUser{
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

	// defer request done
	defer h.doneRequest(r)

	// get user
	usr := updateUser{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.log.Error("Func CreateUser", zap.String("error", err.Error()))
		errDistributor(errs.ErrBadRequestBody, w)
		return
	}

	// creating & transform models.User -> user
	err = h.service.Create(h.updateContext(r.Context(), r), *models.NewUser(usr.ID, usr.Name, usr.Password))
	if err != nil {
		h.log.Error("Func CreateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
	w.Write([]byte("User Created"))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Start func UpdateUser", zap.String("user_id", r.PathValue("user_id")))

	// defer request done
	defer h.doneRequest(r)

	// check path
	// id, err := strconv.Atoi(r.PathValue("user_id"))
	// if err != nil {
	// 	h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
	// 	errDistributor(errs.ErrBadRequestQuery, w)
	// 	return
	// }

	// get user
	usr := getUser{}
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}

	// updating
	err = h.service.Update(h.updateContext(r.Context(), r), models.User{
		Name: usr.Name,
		ID:   usr.ID,
	})
	if err != nil {
		h.log.Error("Func UpdateUser", zap.String("error", err.Error()))
		errDistributor(err, w)
		return
	}
	w.Write([]byte("User Updated"))
}
