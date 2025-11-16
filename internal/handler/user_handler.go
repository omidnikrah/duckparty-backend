package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_dto "github.com/omidnikrah/duckparty-backend/internal/dto/user"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Authenticate(c *gin.Context) {
	var requestBody user_dto.AuthenticateUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	otpErr := h.userService.SendOTP(requestBody.Email, c.Request.Context())
	if otpErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": otpErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Otp sent!",
	})
}

func (h *UserHandler) AuthenticateVerify(c *gin.Context) {
	var requestBody user_dto.AuthenticateUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	user, token, err := h.userService.AuthenticateUser(requestBody.Email, requestBody.OTP, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
