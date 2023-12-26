package handler

import (
	"NewsChanel/internal/helper"
	"NewsChanel/internal/models"
	"NewsChanel/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HandleNotifications(ctx *gin.Context, store *service.NotificationStore) {
	userID, err := helper.GetUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK,
			gin.H{
				"message":       "No notifications found for user:" + userID,
				"notifications": []models.Notification{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}
