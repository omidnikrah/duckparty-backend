package model

import (
	"time"

	"github.com/omidnikrah/duckparty-backend/internal/types"
	"gorm.io/gorm"
)

type Duck struct {
	ID            uint                 `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	DeletedAt     gorm.DeletedAt       `json:"deleted_at,omitempty" gorm:"index"`
	OwnerID       uint                 `json:"owner_id" gorm:"not null;index"`
	Owner         User                 `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name          string               `json:"name" gorm:"not null"`
	X             float64              `json:"x" gorm:"not null"`
	Y             float64              `json:"y" gorm:"not null"`
	Appearance    types.DuckAppearance `json:"appearance" gorm:"serializer:json;type:jsonb;not null"`
	Image         string               `json:"image" gorm:"not null"`
	LikesCount    int64                `json:"likes_count" gorm:"not null;default:0"`
	DislikesCount int64                `json:"dislikes_count" gorm:"not null;default:0"`
	Rank          uint                 `json:"rank" gorm:"not null;default:0"`
}
