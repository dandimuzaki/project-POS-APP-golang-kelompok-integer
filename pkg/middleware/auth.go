package middleware

import (
	"errors"
	"net/http"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareCustom) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			utils.ResponseFailed(c, http.StatusUnauthorized, "missing token", nil)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		userID, err := mw.Usecase.AuthService.ValidateToken(c, token)
		if err != nil {
			utils.ResponseFailed(c, http.StatusUnauthorized, "invalid token", err)
			c.Abort()
			return
		}

		user, err := mw.Usecase.UserService.GetUserByID(c, *userID)
		if err != nil {
			utils.ResponseFailed(c, http.StatusUnauthorized, "user not found", err)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", entity.UserRole(user.Role))
		c.Next()
	}
}

func (mw *MiddlewareCustom) RequirePermission(roles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			utils.ResponseFailed(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		userRole, ok := roleVal.(entity.UserRole)
		if !ok {
			utils.ResponseFailed(c, http.StatusInternalServerError, "invalid role type", nil)
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		utils.ResponseFailed(
			c,
			http.StatusForbidden,
			"forbidden",
			errors.New("insufficient permission"),
		)
		c.Abort()
	}
}
