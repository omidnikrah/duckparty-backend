package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/middleware"
	"github.com/omidnikrah/duckparty-backend/internal/model"
	duckService "github.com/omidnikrah/duckparty-backend/internal/service/duck"
)

type DuckHandler struct {
	duckService *duckService.DuckService
}

func NewDuckHandler(duckService *duckService.DuckService) *DuckHandler {
	return &DuckHandler{
		duckService: duckService,
	}
}

// CreateDuck godoc
// @Summary      Create a new duck
// @Description  Creates a new duck with image, name, email, and appearance data
// @Tags         ducks
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image       formData  file    true   "Duck image file"
// @Param        name        formData  string  true   "Duck name"
// @Param        appearance  formData  string  true   "Duck appearance JSON"
// @Success      200         {object}  duck_dto.DuckResponse  "Created duck"
// @Failure      400         {object}  map[string]string  "Error message"
// @Failure      500         {object}  map[string]string  "Error message"
// @Router       /duck [post]
func (h *DuckHandler) CreateDuck(c *gin.Context) {
	name := c.PostForm("name")
	appearanceJSON := c.PostForm("appearance")

	user, _ := middleware.GetAuthUser(c)

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	if name == "" || appearanceJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, and appearance are required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file: " + err.Error()})
		return
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file content: " + err.Error()})
		return
	}

	req := duckService.CreateDuckRequest{
		Name:           name,
		Email:          user.Email,
		AppearanceJSON: appearanceJSON,
		ImageData:      fileContent,
	}

	newDuck, err := h.duckService.CreateDuck(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDuck)
}

// ReactionToDuck godoc
// @Summary      React to a duck
// @Description  Add a like or dislike reaction to a duck
// @Tags         ducks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        duckId     path      string  true   "Duck ID"
// @Param        reaction   path      string  true   "Reaction type (like or dislike)"  Enums(like, dislike)
// @Success      200        {object}  duck_dto.DuckReactionResponse  "Reaction created"
// @Failure      400        {object}  map[string]string    "Error message"
// @Failure      404        {object}  map[string]string    "Duck not found"
// @Failure      409        {object}  map[string]string    "Duck already reacted"
// @Failure      500        {object}  map[string]string    "Error message"
// @Router       /duck/{duckId}/reaction/{reaction} [put]
func (h *DuckHandler) ReactionToDuck(c *gin.Context) {
	duckId, err := strconv.ParseUint(c.Param("duckId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid duck id"})
		return
	}
	reaction := model.ReactionType(c.Param("reaction"))

	user, _ := middleware.GetAuthUser(c)

	req := duckService.ReactToDuckRequest{DuckID: uint(duckId), UserID: user.UserID, Reaction: reaction}

	duck, err := h.duckService.ReactionToDuck(req)
	if err != nil {
		switch {
		case errors.Is(err, duckService.ErrDuckNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, duckService.ErrDuckAlreadyReacted):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, duck)
}

// GetDucksList godoc
// @Summary      Get list of ducks
// @Description  Returns a list of all ducks ordered by creation date
// @Tags         ducks
// @Accept       json
// @Produce      json
// @Success      200  {array}   duck_dto.DuckResponse  "List of ducks"
// @Failure      500  {object}  map[string]string  "Error message"
// @Router       /ducks [get]
func (h *DuckHandler) GetDucksList(c *gin.Context) {
	ducks, err := h.duckService.GetDucksList()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ducks)
}

// GetUserDucks godoc
// @Summary      Get list of ducks for a specific user
// @Description  Returns a list of all ducks owned by the specified user, ordered by creation date
// @Tags         ducks
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Success      200      {array}   duck_dto.DuckResponse  "List of user's ducks"
// @Failure      400      {object}  map[string]string  "Error message"
// @Failure      500      {object}  map[string]string  "Error message"
// @Router       /user/{userId}/ducks [get]
func (h *DuckHandler) GetUserDucks(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	ducks, err := h.duckService.GetUserDucksList(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ducks)
}

// GetDucksLeaderboard godoc
// @Summary      Get ducks leaderboard
// @Description  Returns the top 100 ducks sorted by rank (highest to lowest)
// @Tags         ducks
// @Accept       json
// @Produce      json
// @Success      200  {array}   duck_dto.DuckResponse  "List of top 100 ducks by rank"
// @Failure      500  {object}  map[string]string  "Error message"
// @Router       /leaderboard [get]
func (h *DuckHandler) GetDucksLeaderboard(c *gin.Context) {
	ducks, err := h.duckService.GetDucksLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ducks)
}

// RemoveDuck godoc
// @Summary      Remove a duck
// @Description  Deletes a duck owned by the authenticated user
// @Tags         ducks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        duckId   path      int  true  "Duck ID"
// @Success      200      {object}  map[string]string  "Success message"
// @Failure      400      {object}  map[string]string  "Invalid duck ID"
// @Failure      404      {object}  map[string]string  "Duck not found"
// @Failure      500      {object}  map[string]string  "Error message"
// @Router       /duck/{duckId} [delete]
func (h *DuckHandler) RemoveDuck(c *gin.Context) {
	authUser, _ := middleware.GetAuthUser(c)
	duckId, err := strconv.ParseUint(c.Param("duckId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid duck id"})
		return
	}

	_, err = h.duckService.RemoveDuck(authUser.UserID, uint(duckId))
	if err != nil {
		switch {
		case errors.Is(err, duckService.ErrDuckNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Duck removed successfully"})
}
