package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", errors.New("message not found")
	}
	return userID, nil
}
