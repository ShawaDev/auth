package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/ShawaDev/auth/pkg/model"
	"github.com/ShawaDev/auth/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	tokenTTL = 12 * time.Hour
	key      = "fakqmrlqmlvsopk,qer"
	salt     = "12,;NMFONQNLKAD"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generateHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, string, error) {
	if s.repo == nil {
		return "", "", errors.New("repository is not initialized")
	}

	user, err := s.repo.GetUser(email, generateHash(password))
	if err != nil {
		return "", "", err
	}

	if user.Id == 0 || user.Role == "" {
		return "", "", errors.New("invalid user data")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(tokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	})

	if key == "" {
		return "", "", errors.New("signing key is not initialized")
	}

	accessToken, err := token.SignedString([]byte(key))

	return accessToken, user.Role, err

}

// Hash function for password
func generateHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
