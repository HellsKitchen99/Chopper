package http

import (
	"chopper/internal/domain"
	"chopper/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *usecase.UserService
}

func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) RegisterRoutes(router gin.IRouter) {
	router.POST("/register", u.UserRegister)
}

func (u *UserHandler) UserRegister(c *gin.Context) {
	ctx := c.Request.Context()
	var userRegisterFromFront domain.UserRegisterFromFront
	if err := c.ShouldBindJSON(&userRegisterFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	if err := u.userService.CreateUser(ctx, userRegisterFromFront); err != nil && errors.Is(err, usecase.ErrUserExists) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user already exists",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.Status(http.StatusCreated)
}
