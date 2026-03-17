package handler

import (
	"net/http"

	"WBTech_L3.1/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type resultResponse struct {
	Result map[string]any `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	ID     string       `json:"id"`
	Status model.Status `json:"status"`
}

func ReturnErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func ReturnResultResponse(ctx *gin.Context, result map[string]any) {
	ctx.JSON(http.StatusOK, resultResponse{result})
}
