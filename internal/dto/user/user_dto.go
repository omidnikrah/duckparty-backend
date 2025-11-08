package user_dto

type AuthenticateUserDTO struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp"`
}
