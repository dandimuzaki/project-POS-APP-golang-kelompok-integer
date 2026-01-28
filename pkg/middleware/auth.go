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

func (mw *MiddlewareCustom) AuthMiddleware() func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			auth := c.Request.Header.Get("Authorization")
			token := strings.TrimSpace(strings.Replace(auth, "Bearer", "", 1))
			userID, err := mw.Usecase.AuthService.ValidateToken(c, token)
			if err != nil {
				utils.ResponseFailed(c, http.StatusUnauthorized, "invalid token", err)
				return
			}

			user, err := mw.Usecase.UserService.GetByID(c, *userID)
			if err != nil {
				utils.ResponseFailed(c, http.StatusUnauthorized, "user not found", err)
				return
			}

			c.Set("user", user)
			c.Next()
		}
	}
}