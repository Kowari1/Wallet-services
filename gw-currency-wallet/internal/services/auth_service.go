package services

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/internal/utils"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   storages.UserStorage
	walletRepo storages.WalletStorage
	jwt        *JWTManager
}

func NewAuthService(userRepo storages.UserStorage, walletRepo storages.WalletStorage, jwt *JWTManager) *AuthService {
	return &AuthService{userRepo: userRepo, walletRepo: walletRepo, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (uuid.UUID, error) {
	exists, err := s.userRepo.ExistUserByUsername(ctx, req.Username)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, ErrUserAlreadyExists
	}

	exists, err = s.userRepo.ExistUserByEmail(ctx, req.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, ErrEmailAlreadyExists
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return uuid.Nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

var (
	ErrUserAlreadyExists  = fmt.Errorf("username already exists")
	ErrEmailAlreadyExists = fmt.Errorf("email already exists")
	ErrInvalidCredentials = fmt.Errorf("invalid username or password")
)
