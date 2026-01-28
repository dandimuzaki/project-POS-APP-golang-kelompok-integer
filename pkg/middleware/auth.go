package middleware

import (
	"net/http"
	"project-POS-APP-golang-integer/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := "2"
		c.Set("ctxID", ID)
		c.Next()
	}
}

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

		user, err := mw.Usecase.UserService.GetByID(c, *userID)
		if err != nil {
			utils.ResponseFailed(c, http.StatusUnauthorized, "user not found", err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}