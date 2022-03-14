package http

import (
	"go-midtrans/domain"
	"go-midtrans/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(r *gin.RouterGroup, au domain.AuthUsecase) {
	handler := UserHandler{authUsecase: au}
	api := r.Group("/auth")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}
}

type RegisInput struct {
	Name     string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type LoginInput struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var input RegisInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &utils.Response{Message: err.Error()})
		return
	}

	user, err := h.authUsecase.Register(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(utils.GetErrCode(err), &utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &utils.Response{Data: AuthResponse(user, ""), Message: "account successfully registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &utils.Response{Message: err.Error()})
		return
	}

	user, err := h.authUsecase.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(utils.GetErrCode(err), &utils.Response{Message: err.Error()})
		return
	}

	claim := jwt.MapClaims{}
	claim["email"] = user.Email

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &utils.Response{Data: AuthResponse(user, signedToken), Message: "successfully login"})
}
