package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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
