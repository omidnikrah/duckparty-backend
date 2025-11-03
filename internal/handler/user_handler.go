package handler

import (
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
