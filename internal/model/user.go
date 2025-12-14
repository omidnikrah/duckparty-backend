package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       *string `json:"email" gorm:"unique"`
	DisplayName *string `json:"display_name"`
}
