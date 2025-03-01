package utils

import (
	"github.com/gin-gonic/gin"
)

// StatusType represents the API response status.
type StatusType string

func CreateJSONResponse(status StatusType, message string) gin.H {

	response := gin.H{
		"result": gin.H{
			"status":  status,
			"message": message,
		},
	}

	return response
}
