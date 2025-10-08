package model

import (
	"github.com/omidnikrah/duckparty-backend/internal/types"
	"gorm.io/gorm"
)

type Duck struct {
	gorm.Model
	OwnerID    uint                 `json:"owner_id" gorm:"not null;index"`
	Owner      User                 `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name       string               `json:"name" gorm:"not null"`
	X          float64              `json:"x" gorm:"not null"`
	Y          float64              `json:"y" gorm:"not null"`
	Appearance types.DuckAppearance `json:"appearance" gorm:"serializer:json;type:jsonb;not null"`
	LikesCount int64                `json:"likes_count" gorm:"not null;default:0"`
}
