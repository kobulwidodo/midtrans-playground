package http

import (
	"go-midtrans/domain"
	"go-midtrans/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type middleware struct {
	authUsecase domain.AuthUsecase
}

func NewMiddleware(au domain.AuthUsecase) gin.HandlerFunc {
	return (&middleware{authUsecase: au}).Handle()
}

func (m *middleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &utils.Response{Message: "wrong token type"})
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		user, err := m.authUsecase.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &utils.Response{Message: "token invalid"})
			return
		}

		c.Set("userLoggedin", user)
	}
}
