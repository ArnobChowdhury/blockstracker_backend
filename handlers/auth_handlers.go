package handlers

import (
	"fmt"
	"net/http"

	responsemsg "blockstracker_backend/constants"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		var errors []string

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field %s is invalid: %s", e.Field(), e.Tag()))
		}

		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responsemsg.UserCreationSuccess})
}
