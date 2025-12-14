package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_dto "github.com/omidnikrah/duckparty-backend/internal/dto/user"
	"github.com/omidnikrah/duckparty-backend/internal/middleware"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Authenticate godoc
// @Summary      Send OTP to user email
// @Description  Sends a one-time password (OTP) to the user's email address for authentication
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_dto.AuthenticateUserDTO  true  "Email address"
// @Success      200      {object}  map[string]string            "Success message"
// @Failure      400      {object}  map[string]string            "Error message"
// @Router       /auth [post]
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

// AuthenticateVerify godoc
// @Summary      Verify OTP and authenticate user
// @Description  Verifies the OTP code and returns user information along with JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_dto.AuthenticateUserDTO  true  "Email and OTP code"
// @Success      200      {object}  user_dto.AuthenticateResponse  "User and token"
// @Failure      400      {object}  map[string]string            "Error message"
// @Router       /auth/verify [post]
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

// UpdateName godoc
// @Summary      Update user display name
// @Description  Updates the display name of the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      user_dto.UpdateNameDTO  true  "New display name"
// @Success      200      {object}  user_dto.UserInfoResponse  "Updated user"
// @Failure      400      {object}  map[string]string            "Error message"
// @Router       /user/change-name [put]
func (h *UserHandler) UpdateName(c *gin.Context) {
	var requestBody user_dto.UpdateNameDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	authUser, _ := middleware.GetAuthUser(c)

	updatedUser, err := h.userService.UpdateName(requestBody.Name, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

// SetEmail godoc
// @Summary      Send OTP to new email address
// @Description  Sends an OTP to the new email address for verification. Use /user/verify-email to verify and set the email.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      user_dto.SetEmailDTO  true  "Email address"
// @Success      200      {object}  map[string]string            "Success message"
// @Failure      400      {object}  map[string]string            "Error message"
// @Router       /user/set-email [post]
func (h *UserHandler) SetEmail(c *gin.Context) {
	var requestBody user_dto.SetEmailDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	authUser, _ := middleware.GetAuthUser(c)

	if err := h.userService.SetEmail(requestBody.Email, authUser.UserID, c.Request.Context()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email!",
	})
}

// VerifyEmailChange godoc
// @Summary      Verify email change with OTP
// @Description  Verifies the OTP code and sets the email address for the authenticated user. Returns updated user and new token.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      user_dto.AuthenticateUserDTO  true  "Email and OTP code"
// @Success      200      {object}  user_dto.AuthenticateResponse  "Updated user and new token"
// @Failure      400      {object}  map[string]string              "Error message"
// @Router       /user/verify-set-email [post]
func (h *UserHandler) VerifySetEmail(c *gin.Context) {
	var requestBody user_dto.AuthenticateUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	authUser, _ := middleware.GetAuthUser(c)

	updatedUser, token, err := h.userService.VerifySetEmail(requestBody.Email, requestBody.OTP, authUser.UserID, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  updatedUser,
		"token": token,
	})
}

// CreateAnonymousUser godoc
// @Summary      Create anonymous user and get token
// @Description  Creates an anonymous user with a display name and returns a JWT token for immediate use
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_dto.CreateAnonymousUserDTO  true  "User display name"
// @Success      200      {object}  user_dto.AuthenticateResponse      "User and token"
// @Failure      400      {object}  map[string]string                  "Error message"
// @Router       /auth/anonymous [post]
func (h *UserHandler) CreateAnonymousUser(c *gin.Context) {
	var requestBody user_dto.CreateAnonymousUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.Error(err)
		return
	}

	user, token, err := h.userService.CreateAnonymousUser(requestBody.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// GetMeUser godoc
// @Summary      Get current user information
// @Description  Returns the information of the currently authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  user_dto.UserInfoResponse  "User information"
// @Failure      400      {object}  map[string]string            "Error message"
// @Router       /user [get]
func (h *UserHandler) GetMeUser(c *gin.Context) {
	authUser, _ := middleware.GetAuthUser(c)

	meUser, err := h.userService.GetUser(authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": meUser,
	})
}
