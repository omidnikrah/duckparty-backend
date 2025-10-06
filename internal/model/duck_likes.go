package model

import "time"

type DuckLikes struct {
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	DuckID    uint      `json:"duck_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Duck      Duck      `json:"duck" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}
