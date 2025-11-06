package utils

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

func JSONError(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, ErrorResponse{Error: msg})
}

func JSONSuccess(c *gin.Context, code int, data any) {
	c.JSON(code, SuccessResponse{Data: data})
}
