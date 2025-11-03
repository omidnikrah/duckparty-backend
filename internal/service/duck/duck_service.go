package duckService

import (
	"encoding/json"

	"github.com/omidnikrah/duckparty-backend/internal/model"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	"github.com/omidnikrah/duckparty-backend/internal/types"
	"gorm.io/gorm"
)

type DuckService struct {
	db          *gorm.DB
	userService *userService.UserService
	storage     *storage.S3Storage
}

func NewService(db *gorm.DB, userService *userService.UserService, s3Storage *storage.S3Storage) *DuckService {
	return &DuckService{
		db:          db,
		userService: userService,
		storage:     s3Storage,
	}
}

type CreateDuckRequest struct {
	Name           string
	Email          string
	AppearanceJSON string
	ImageData      []byte
}

func (s *DuckService) CreateDuck(req CreateDuckRequest) (*model.Duck, error) {
	var appearance types.DuckAppearance
	if err := json.Unmarshal([]byte(req.AppearanceJSON), &appearance); err != nil {
		return nil, err
	}

	imageURL, err := s.storage.UploadFile(req.ImageData, req.Name)
	if err != nil {
		return nil, err
	}

	var newDuck model.Duck

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.userService.GetOrCreateUserByEmail(req.Email, tx)
		if err != nil {
			return err
		}

		newDuck = model.Duck{
			OwnerID:    user.ID,
			Name:       req.Name,
			Appearance: appearance,
			Image:      imageURL,
		}

		if err := tx.Create(&newDuck).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newDuck, nil
}
