package duckService

import (
	"encoding/json"
	"errors"

	"github.com/omidnikrah/duckparty-backend/internal/model"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	"github.com/omidnikrah/duckparty-backend/internal/types"
	"gorm.io/gorm"
)

var (
	ErrDuckNotFound       = errors.New("duck not found")
	ErrDuckAlreadyReacted = errors.New("duck already reacted")
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

type ReactToDuckRequest struct {
	DuckID   uint
	UserID   uint
	Reaction model.ReactionType
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

		if err := tx.Preload("Owner").First(&newDuck, newDuck.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newDuck, nil
}

func (s *DuckService) ReactionToDuck(req ReactToDuckRequest) (*model.DuckReactions, error) {
	var (
		reaction model.DuckReactions
		duck     model.Duck
	)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var existingReaction model.DuckReactions

		if err := tx.Preload("Owner").First(&duck, req.DuckID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrDuckNotFound
			}
			return err
		}

		if err := tx.Where("duck_id = ? AND user_id = ?", req.DuckID, req.UserID).First(&existingReaction).Error; err == nil {
			if existingReaction.Reaction == req.Reaction {
				return ErrDuckAlreadyReacted
			}

			if err := tx.Delete(&existingReaction).Error; err != nil {
				return err
			}

			updateReactionCounts(&duck, existingReaction.Reaction, -1)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		reaction = model.DuckReactions{
			DuckID:   req.DuckID,
			UserID:   req.UserID,
			Reaction: req.Reaction,
		}

		if err := tx.Create(&reaction).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuckAlreadyReacted
			}
			return err
		}

		updateReactionCounts(&duck, req.Reaction, 1)

		if err := tx.Save(&duck).Error; err != nil {
			return err
		}

		reaction.Duck = duck

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (s *DuckService) GetDucksList() (*[]model.Duck, error) {
	ducks := []model.Duck{}
	if err := s.db.Preload("Owner").Order("created_at DESC").Find(&ducks).Error; err != nil {
		return nil, err
	}

	return &ducks, nil
}

func updateReactionCounts(duck *model.Duck, reaction model.ReactionType, delta int64) {
	switch reaction {
	case model.ReactionLike:
		duck.LikesCount = clampNonNegative(duck.LikesCount + delta)
	case model.ReactionDislike:
		duck.DislikesCount = clampNonNegative(duck.DislikesCount + delta)
	}
}

func clampNonNegative(value int64) int64 {
	if value < 0 {
		return 0
	}
	return value
}
