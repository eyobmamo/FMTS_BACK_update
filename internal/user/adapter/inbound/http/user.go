package user_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	user "FMTS/internal/user/application"
	port "FMTS/internal/user/port/inbound/user"
	context "FMTS/pkg/context"
	"FMTS/pkg/utils"

	utility "FMTS/utils"
)

type UserHandler struct {
	userService user.UserService
	logger      utils.Logger
}

func NewUserHandler(service user.UserService, logger utils.Logger) port.UserPortHandler {
	return &UserHandler{
		userService: service,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[CreateUser] failed to decode request: %v", err)
		utility.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Warnf("[CreateUser] validation failed: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	ctx := context.ExtractUserContext(r)
	fmt.Println("=====================================")
	fmt.Printf("user id: %v", ctx.UserID)
	fmt.Println("=====================================")

	if ctx.UserID == "" {
		utility.SendErrorResponse(w, "unauthorized: user context missing", http.StatusUnauthorized, nil)
		return
	}
	userCreated, err := h.userService.CreateUser(req, ctx.UserID)
	if err != nil {
		h.logger.Errorf("[CreateUser] service error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, userCreated, "User Created Sussesfully")

}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userFetched, err := h.userService.GetUserByID(id)
	if err != nil {
		h.logger.Errorf("[GetUserByID] failed: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}
	utility.WriteSuccessResponse(w, userFetched, "User fetched successfully")
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers()
	if err != nil {
		h.logger.Errorf("[ListUsers] failed: %v", err)
		utility.SendErrorResponse(w, "failed to list users", http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, users, "Users retrieved")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateUser] decode error: %v", err)
		utility.SendErrorResponse(w, "invalid input", http.StatusBadRequest, nil)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Warnf("[UpdateUser] validation error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	updated, err := h.userService.UpdateUser(id, req)
	if err != nil {
		h.logger.Errorf("[UpdateUser] update error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, updated, "User updated")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		h.logger.Errorf("[DeleteUser] delete error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, "given user : 234234", "User deleted sucessfuly")
}
