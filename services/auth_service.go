package services

import (
	"errors"
	"time"

	"auth-service/config"
	"auth-service/models"
	"auth-service/repositories"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	repo *repositories.UserRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repositories.NewUserRepo(),
	}
}

func (s *AuthService) Register(email, password, fullName string) (*models.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("email already used")
	}
	user := &models.User{
		Email:    email,
		FullName: fullName,
		Role:     "user",
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !u.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) GenerateJWT(u *models.User) (string, error) {
	ttl := time.Duration(config.JWTTTLMin) * time.Minute
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"role":  u.Role,
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
		"email": u.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}
