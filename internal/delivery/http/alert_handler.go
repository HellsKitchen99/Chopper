package http

import (
	"chopper/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlertHandler struct {
	alertService *usecase.AlertService
}

func NewAlertHandler(alertService *usecase.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

func (a *AlertHandler) RegisterRoutes(r gin.IRouter) {
	r.GET("/get", a.GetLastSevenDaysAlert)
}

func (a *AlertHandler) GetLastSevenDaysAlert(c *gin.Context) {
	ctx := c.Request.Context()
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad credentials",
		})
		return
	}
	userId, ok := uid.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad credentials",
		})
		return
	}
	allertMessage, err := a.alertService.GetLastSevenDays(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": allertMessage,
	})
}
