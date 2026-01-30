package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service usecase.UserService
	Logger *zap.Logger
	Config utils.Configuration
}

func NewUserHandler(service usecase.UserService, log *zap.Logger, config utils.Configuration) UserHandler {
	return UserHandler{
		service: service,
		Logger: log,
		Config: config,
	}
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	// Construct DTO
	req := request.UserFilterRequest{
		Role: c.Query("role"),
		Name: c.Query("name"),
		Email: c.Query("email"),
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	result, err := h.service.GetUserList(c, req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadGateway, "get user list failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get user list success", result)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	role := c.Value("user_role").(entity.UserRole)
	if role != entity.RoleSuperAdmin && req.Role == string(entity.RoleSuperAdmin) {
		utils.ResponseFailed(c, http.StatusForbidden, "create user failed", err)
		return
	}

	res, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "create user failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "create user success", res)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	role := c.Value("user_role").(entity.UserRole)
	if role != entity.RoleSuperAdmin && req.Role == string(entity.RoleSuperAdmin) {
		utils.ResponseFailed(c, http.StatusForbidden, "update user failed", err)
		return
	}

	err = h.service.UpdateRole(c.Request.Context(), req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "update user failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "update user success", nil)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := h.service.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "delete user failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "delete user success", nil)
}