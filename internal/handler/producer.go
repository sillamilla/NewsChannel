package handler

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SendNewsHandler(producer sarama.SyncProducer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h.srv.SendNews(producer, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}
