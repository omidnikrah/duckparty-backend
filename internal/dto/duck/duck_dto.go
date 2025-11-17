package duck_dto

import (
	"time"

	"github.com/omidnikrah/duckparty-backend/internal/model"
	"github.com/omidnikrah/duckparty-backend/internal/types"
)

type CreateDuckDTO struct {
	Name       string               `json:"name" binding:"required"`
	Image      string               `json:"image" binding:"required"`
	Email      string               `json:"email" binding:"required,email"`
	Appearance types.DuckAppearance `json:"appearance" binding:"required"`
}

type ReactToDuckDTO struct {
	DuckId   uint               `json:"duck_id" binding:"required"`
	Reaction model.ReactionType `json:"reaction" binding:"required"`
}

type DuckResponse struct {
	ID            uint                 `json:"id" example:"1"`
	CreatedAt     time.Time            `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time            `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	OwnerID       uint                 `json:"owner_id" example:"1"`
	Owner         DuckUserResponse     `json:"owner"`
	Name          string               `json:"name" example:"Ducky"`
	X             float64              `json:"x" example:"100.5"`
	Y             float64              `json:"y" example:"200.5"`
	Appearance    types.DuckAppearance `json:"appearance"`
	Image         string               `json:"image" example:"https://example.com/image.jpg"`
	LikesCount    int64                `json:"likes_count" example:"10"`
	DislikesCount int64                `json:"dislikes_count" example:"2"`
	Rank          uint                 `json:"rank" example:"1"`
} // @name DuckResponse

type DuckUserResponse struct {
	ID          uint      `json:"id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	Email       string    `json:"email" example:"user@example.com"`
	DisplayName string    `json:"display_name" example:"John Doe"`
} // @name DuckUserResponse

type DuckReactionResponse struct {
	UserID    uint               `json:"user_id" example:"1"`
	DuckID    uint               `json:"duck_id" example:"1"`
	Reaction  model.ReactionType `json:"reaction" example:"like" enums:"like,dislike"`
	User      DuckUserResponse   `json:"user"`
	Duck      DuckResponse       `json:"duck"`
	CreatedAt time.Time          `json:"created_at" example:"2024-01-01T00:00:00Z"`
} // @name DuckReactionResponse
