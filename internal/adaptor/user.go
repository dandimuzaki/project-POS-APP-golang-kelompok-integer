package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto"
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
	ctx := c.Request.Context()

	role := c.Query("role")

	// Construct DTO
	req := dto.UserFilterRequest{
		Role: role,
	}

	result, pagination, err := h.service.UserService.GetUserList(ctx, req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadGateway, "", nil)
		return
	}

	utils.ResponsePagination(c, http.StatusOK, "success get data", result, pagination)
}