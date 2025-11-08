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

func (h *DuckHandler) CreateDuck(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	appearanceJSON := c.PostForm("appearance")

	if name == "" || email == "" || appearanceJSON == "" {
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
		Email:          email,
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
