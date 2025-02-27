package handlers

import (
	"errors"
	"log"
	"net/http"

	responsemsg "blockstracker_backend/constants"
	"blockstracker_backend/internal/repositories"

	"blockstracker_backend/internal/validators"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
}

func NewAuthHandler(userRepo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			message := validators.GetCustomMessage(validationErrors[0], req)
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": responsemsg.MalformedRequest})
		return
	}

	hashedPassword, pwHashingError := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if pwHashingError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": responsemsg.InternalServerError})
		return
	}

	hashedPasswordString := string(hashedPassword)

	user := models.User{
		Email:    req.Email,
		Password: &hashedPasswordString,
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": responsemsg.UserCreationFailed})
			return
		} else {
			log.Printf("%s: %v", responsemsg.UnexpectedErrorDuringUserCreation, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": responsemsg.InternalServerError})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": responsemsg.UserCreationSuccess})
}
