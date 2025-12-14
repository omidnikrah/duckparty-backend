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
	"github.com/omidnikrah/duckparty-backend/internal/templates"
	"github.com/omidnikrah/duckparty-backend/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	authTokenTTL    = 30 * 24 * time.Hour // 30 days
	authRedisTTL    = 2 * time.Minute
	otpEmailSubject = "DuckParty OTP Code"
)

type UserService struct {
	db           *gorm.DB
	rdb          *redis.Client
	config       *config.Config
	resendClient *resend.Client
}

func NewService(db *gorm.DB, rdb *redis.Client, resendClient *resend.Client, config *config.Config) *UserService {
	return &UserService{db: db, rdb: rdb, resendClient: resendClient, config: config}
}

func (s *UserService) CreateUserByName(name string, tx *gorm.DB) (*model.User, error) {
	db := s.db
	if tx != nil {
		db = tx
	}

	newUser := model.User{DisplayName: &name}
	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *UserService) CreateAnonymousUser(name string) (*model.User, string, error) {
	user, err := s.CreateUserByName(name, nil)
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) GetOrCreateUserByEmail(email string, tx *gorm.DB) (*model.User, error) {
	db := s.db
	if tx != nil {
		db = tx
	}

	newUser := model.User{Email: &email}
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

	if err := s.rdb.Set(ctx, key, otpCode, authRedisTTL).Err(); err != nil {
		return err
	}

	if err := s.sendOtpEmail(email, otpCode); err != nil {
		return err
	}

	return nil
}

func (s *UserService) verifyOTP(email string, otp string, ctx context.Context) error {
	key := fmt.Sprintf("otp:user:%s", email)

	storedOtp, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("otp expired or invalid")
		}
		return err
	}

	if storedOtp != otp {
		return errors.New("invalid otp")
	}

	if err := s.rdb.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (s *UserService) GenerateToken(user *model.User) (string, error) {
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
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) AuthenticateUser(email string, otp string, ctx context.Context) (*model.User, string, error) {
	if err := s.verifyOTP(email, otp, ctx); err != nil {
		return nil, "", err
	}

	user, err := s.GetOrCreateUserByEmail(email, nil)
	if err != nil {
		return nil, "", err
	}

	tokenString, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

func (s *UserService) sendOtpEmail(email string, otpCode int) error {
	htmlBody, err := templates.GenerateOTPEmailHTML(otpCode)
	if err != nil {
		return fmt.Errorf("failed to render HTML email: %w", err)
	}

	textBody, err := templates.GenerateOTPEmailText(otpCode)
	if err != nil {
		return fmt.Errorf("failed to render text email: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    s.config.AuthSenderEmail,
		To:      []string{email},
		Subject: otpEmailSubject,
		Html:    htmlBody,
		Text:    textBody,
	}

	_, err = s.resendClient.Emails.Send(params)
	if err != nil {
		return errors.New("unable to send email. Please try again later")
	}

	return nil
}

func (s *UserService) UpdateName(name string, userId uint) (model.User, error) {
	var user model.User
	err := s.db.Model(&model.User{}).Where("id = ?", userId).Update("display_name", name).Error
	if err != nil {
		return model.User{}, err
	}

	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) SetEmail(email string, userId uint, ctx context.Context) error {
	pendingEmailKey := fmt.Sprintf("set_email:user:%s", email)

	var count int64
	if err := s.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email already exists")
	}

	if err := s.rdb.Set(ctx, pendingEmailKey, userId, authRedisTTL).Err(); err != nil {
		return err
	}

	if err := s.SendOTP(email, ctx); err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifySetEmail(email string, otp string, userId uint, ctx context.Context) (*model.User, string, error) {
	if err := s.verifyOTP(email, otp, ctx); err != nil {
		return nil, "", err
	}

	pendingEmailKey := fmt.Sprintf("set_email:user:%s", email)
	pendingUserIdStr, err := s.rdb.Get(ctx, pendingEmailKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, "", errors.New("no pending email change found")
		}
		return nil, "", err
	}

	pendingUserId, err := strconv.ParseUint(pendingUserIdStr, 10, 64)
	if err != nil {
		return nil, "", errors.New("invalid pending email change")
	}

	if uint(pendingUserId) != userId {
		return nil, "", errors.New("email does not match pending email change for this user")
	}

	if err := s.db.Model(&model.User{}).Where("id = ?", userId).Update("email", email).Error; err != nil {
		return nil, "", err
	}

	user := &model.User{}
	if err := s.db.Where("id = ?", userId).First(user).Error; err != nil {
		return nil, "", err
	}

	tokenString, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	if err := s.rdb.Del(ctx, pendingEmailKey).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return nil, "", err
	}

	return user, tokenString, nil
}

func (s *UserService) GetUser(userId uint) (*model.User, error) {
	user := &model.User{}

	err := s.db.Where("id = ?", userId).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
