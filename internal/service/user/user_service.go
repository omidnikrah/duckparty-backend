package userService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/model"
	"github.com/omidnikrah/duckparty-backend/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const authTokenTTL = 30 * 24 * time.Hour // 30 days

type UserService struct {
	db     *gorm.DB
	rdb    *redis.Client
	config *config.Config
}

func NewService(db *gorm.DB, rdb *redis.Client, config *config.Config) *UserService {
	return &UserService{db: db, rdb: rdb, config: config}
}

func (s *UserService) GetOrCreateUserByEmail(email string, tx *gorm.DB) (*model.User, error) {
	db := s.db
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

func (s *UserService) SendOTP(email string, ctx context.Context) error {
	otpCode := utils.GenerateRandomNumber(5)
	key := fmt.Sprintf("otp:user:%s", email)

	if err := s.rdb.Set(ctx, key, otpCode, 2*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AuthenticateUser(email string, otp string, ctx context.Context) (*model.User, string, error) {
	key := fmt.Sprintf("otp:user:%s", email)

	storedOtp, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, "", errors.New("otp expired or invalid")
		}
		return nil, "", err
	}

	if storedOtp != otp {
		return nil, "", errors.New("invalid otp")
	}

	user, err := s.GetOrCreateUserByEmail(email, nil)
	if err != nil {
		return nil, "", err
	}

	if err := s.rdb.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return nil, "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   strconv.FormatUint(uint64(user.ID), 10),
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(authTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
