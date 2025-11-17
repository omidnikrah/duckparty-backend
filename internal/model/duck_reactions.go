package model

import "time"

type ReactionType string // @name ReactionType

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
)

type DuckReactions struct {
	UserID    uint         `json:"user_id" gorm:"not null;primaryKey;uniqueIndex:idx_duck_user"`
	DuckID    uint         `json:"duck_id" gorm:"not null;primaryKey;uniqueIndex:idx_duck_user"`
	Reaction  ReactionType `json:"reaction" gorm:"type:text;not null;default:'like'"`
	User      User         `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Duck      Duck         `json:"duck" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time    `gorm:"not null;default:now()"`
}
