package duck_dto

import "github.com/omidnikrah/duckparty-backend/internal/types"

type CreateDuckDTO struct {
	Name       string               `json:"name" binding:"required"`
	Image      string               `json:"image" binding:"required"`
	Email      string               `json:"email" binding:"required,email"`
	Appearance types.DuckAppearance `json:"appearance" binding:"required"`
}
