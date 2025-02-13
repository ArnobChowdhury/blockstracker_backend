package handlers

import (
	"net/http"

	responsemsg "blockstracker_backend/constants"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": responsemsg.TaskCreatedSuccess})
}

func UpdateTask(c *gin.Context) {
	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": responsemsg.TaskUpdatedSuccess})
}

func UpdateRepetitiveTask(c *gin.Context) {
	var req models.UpdateRepetitiveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": responsemsg.TaskUpdatedSuccess})
}
