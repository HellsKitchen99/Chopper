package http

import (
	"chopper/internal/domain"
	"chopper/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	dailyNotesService *usecase.DailyNotesService
}

func NewNoteHandler(dailyNotesService *usecase.DailyNotesService) *NoteHandler {
	return &NoteHandler{
		dailyNotesService: dailyNotesService,
	}
}

func (n *NoteHandler) RegisterRoutes(protected gin.IRouter) {
	protected.POST("/new", n.CreateNote)
}

func (n *NoteHandler) CreateNote(c *gin.Context) {
	var dailyNoteFromFront domain.DailyNoteFromFront
	if err := c.ShouldBindJSON(&dailyNoteFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	ctx := c.Request.Context()
	uId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	userId, ok := uId.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "bad token",
		})
		return
	}
	err := n.dailyNotesService.CreateNote(ctx, userId, dailyNoteFromFront)
	if err != nil {
		if errors.Is(err, usecase.ErrNoteAlreadyExists) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		if errors.Is(err, usecase.ErrWrongMoodValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong mood value",
			})
			return
		}
		if errors.Is(err, usecase.ErrWrongSleepHourValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong sleep hours value",
			})
			return
		}
		if errors.Is(err, usecase.ErrWrongLoadValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong load value",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.Status(http.StatusCreated)
}

func (n *NoteHandler) ChangeMood(c *gin.Context) {
	var changeMoodFromFront domain.ChangeMoodFromFront
	if err := c.ShouldBindJSON(&changeMoodFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request body",
		})
		return
	}
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
	changeMessage, err := n.dailyNotesService.ChangeMood(ctx, userId, changeMoodFromFront.Date, changeMoodFromFront.Mood)
	if err != nil {
		if errors.Is(err, usecase.ErrNoteNotExists) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "note not exists",
			})
			return
		}
		if errors.Is(err, usecase.ErrWrongMoodValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong mood value",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"answer": changeMessage,
	})
}

func (n *NoteHandler) ChangeSleepHours(c *gin.Context) {
	var changeSleepHoursFromFront domain.ChangeSleepHoursFromFront
	if err := c.ShouldBindJSON(&changeSleepHoursFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request body",
		})
		return
	}
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
	changeMessage, err := n.dailyNotesService.ChangeSleepHours(ctx, userId, changeSleepHoursFromFront.Date, changeSleepHoursFromFront.SleepHours)
	if err != nil {
		if errors.Is(err, usecase.ErrNoteNotExists) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "note not exists",
			})
			return
		}
		if errors.Is(err, usecase.ErrWrongSleepHourValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong sleep hours value",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"answer": changeMessage,
	})
}

func (n *NoteHandler) ChangeLoad(c *gin.Context) {
	var changeLoadFromFront domain.ChangeLoadFromFront
	if err := c.ShouldBindJSON(&changeLoadFromFront); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request body",
		})
		return
	}
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
	changeMessage, err := n.dailyNotesService.ChangeLoad(ctx, userId, changeLoadFromFront.Date, changeLoadFromFront.Load)
	if err != nil {
		if errors.Is(err, usecase.ErrNoteNotExists) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "note not exists",
			})
			return
		}
		if errors.Is(err, usecase.ErrWrongLoadValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "wrong load value",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"answer": changeMessage,
	})
}
