package user_dto

import "time"

type AuthenticateUserDTO struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP   string `json:"otp" example:"123456"`
} // @name AuthenticateRequest

type AuthenticateResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
} // @name AuthenticateResponse

type UserResponse struct {
	ID          uint      `json:"id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	Email       string    `json:"email" example:"user@example.com"`
	DisplayName string    `json:"display_name" example:"John Doe"`
} // @name UserResponse

type UpdateNameDTO struct {
	Name string `json:"name" binding:"required"`
} // @name UpdateNameRequest
