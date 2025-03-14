package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email      string         `gorm:"not null;unique" json:"email"`
	Password   *string        `gorm:"type:varchar" json:"-"`        // Nullable, hidden in JSON
	Provider   *string        `gorm:"type:varchar" json:"provider"` // Nullable
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	ModifiedAt time.Time      `gorm:"autoUpdateTime" json:"modifiedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,strongpassword" example:"Strongpassword123"`
}

type EmailSignInRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required" example:"Strongpassword123"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// Response example structs for swagger
type GenericSuccessResponse struct {
	Result SuccessResult `json:"result"`
}

type GenericErrorResponse struct {
	Result ErrorResult `json:"result"`
}

type SignInSuccessResponse struct {
	Result SignInSuccessResult `json:"result"`
}

type SignInSuccessResult struct {
	Status  string        `json:"status" example:"Success"`
	Message string        `json:"message" example:"Success message"`
	Data    TokenResponse `json:"data"`
}

type SuccessResult struct {
	Status  string `json:"status" example:"Success"`
	Message string `json:"message" example:"Success message"`
}

type ErrorResult struct {
	Status  string `json:"status" example:"Error"`
	Message string `json:"message" example:"Error message"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
