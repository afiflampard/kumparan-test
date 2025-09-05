package infra

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func JSONSuccess(c *app.RequestContext, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c *app.RequestContext, statusCode int, message string, err error) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}
