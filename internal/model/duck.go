package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Duck struct {
	gorm.Model
	OwnerID     uint           `json:"owner_id" gorm:"not null;index"`
	Owner       User           `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string         `json:"name" gorm:"not null"`
	X           float64        `json:"x" gorm:"not null"`
	Y           float64        `json:"y" gorm:"not null"`
	Accessories datatypes.JSON `gorm:"type:jsonb;not null"`
	LikesCount  int64          `gorm:"not null;default:0"`
}
