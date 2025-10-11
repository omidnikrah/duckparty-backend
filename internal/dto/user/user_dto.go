package user_dto

type CreateUserDTO struct {
	Email string `json:"email" binding:"required,email"`
}
