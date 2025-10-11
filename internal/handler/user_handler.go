package handler

import (
	"github.com/omidnikrah/duckparty-backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetOrCreateUserByEmail(email string, tx *gorm.DB) (*model.User, error) {
	db := h.db
	if tx != nil {
		db = tx
	}

	newUser := model.User{Email: email}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&newUser).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
