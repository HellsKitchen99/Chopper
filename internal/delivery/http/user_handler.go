package http

import (
	"chopper/internal/domain"
	"chopper/internal/usecase"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *usecase.UserService
}

func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) RegisterRoutes(public gin.IRouter, protected gin.IRouter) {
	public.POST("/register", u.UserRegister)
	public.POST("/login", u.UserLogin)
	protected.GET("/me", u.WhoAmI)
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

func (u *UserHandler) UserLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var userLoginFromFront domain.UserLoginFromFront
	if err := c.ShouldBindJSON(&userLoginFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	token, err := u.userService.CheckUserInDatabase(ctx, userLoginFromFront)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotExist) || errors.Is(err, usecase.ErrWrongPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.Header("Authorization", fmt.Sprintf("Bearer %v", token))
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (u *UserHandler) WhoAmI(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	userId, ok := id.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	userUsername, ok := username.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	ctx := c.Request.Context()
	user, err := u.userService.GetIdUsernameRole(ctx, userId, userUsername)
	if err != nil && errors.Is(err, usecase.ErrUserNotExist) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
