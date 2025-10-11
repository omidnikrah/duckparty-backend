package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/omidnikrah/duckparty-backend/internal/dto/duck"
	"github.com/omidnikrah/duckparty-backend/internal/model"
	"gorm.io/gorm"
)

type DuckHandler struct {
	db          *gorm.DB
	userHandler *UserHandler
}

func NewDuckHandler(db *gorm.DB, userHandler *UserHandler) *DuckHandler {
	return &DuckHandler{
		db:          db,
		userHandler: userHandler,
	}
}

func (h *DuckHandler) CreateDuck(c *gin.Context) {
	var duck dto.CreateDuckDTO
	if err := c.ShouldBindJSON(&duck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newDuck model.Duck

	err := h.db.Transaction(func(tx *gorm.DB) error {
		user, err := h.userHandler.GetOrCreateUserByEmail(duck.Email, tx)
		if err != nil {
			return err
		}

		newDuck = model.Duck{
			OwnerID:    user.ID,
			Name:       duck.Name,
			Appearance: duck.Appearance,
		}

		if err := tx.Create(&newDuck).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDuck)
}
