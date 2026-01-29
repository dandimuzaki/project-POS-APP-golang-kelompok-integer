package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewUserHandler(service *usecase.Usecase, log *zap.Logger, config utils.Configuration) UserHandler {
	return UserHandler{
		service: service,
		Logger: log,
		Config: config,
	}
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	role := c.Query("role")
	name := c.Query("name")
	email := c.Query("email")

	// Construct DTO
	req := request.UserFilterRequest{
		Role: role,
		Name: name,
		Email: email,
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	result, err := h.service.UserService.GetUserList(c, req)
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



	res, err := h.service.UserService.CreateUser(c.Request.Context(), req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "create user failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "create user success", res)
}