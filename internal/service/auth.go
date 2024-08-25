package service

import (
	"autoshop/internal/domain/models"
	"autoshop/internal/service/dto"
	"autoshop/internal/storage"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/hash"
	"autoshop/pkg/jwt_helper"
	"autoshop/pkg/logging"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullName"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=10"`
}

type AuthService interface {
	Register(input RegisterInput) (*dto.TokenResponse, error)
	GetUser(userId uuid.UUID) (*dto.GetUserResponse, error)
	Login(input LoginInput) (*dto.TokenResponse, error)
	Reauth(refresh string) (string, error)
}

type JwtAuthService struct {
	userStore storage.UserStorage
}

func NewJwtAuthService(userStore storage.UserStorage) *JwtAuthService {
	return &JwtAuthService{userStore: userStore}
}

func (s *JwtAuthService) Register(input RegisterInput) (*dto.TokenResponse, error) {
	const op = "JwtAuthService.Register"
	log := logging.CreateLoggerWithOp(op)

	hasher := hash.BcryptHash{}
	passwordHash := hash.HashPassword(hasher, input.Password)

	dbUser := models.User{
		Id:           uuid.New(),
		Email:        input.Email,
		FullName:     input.FullName,
		PasswordHash: passwordHash}

	if err := s.userStore.SaveUser(&dbUser); err != nil {
		if errors.Is(err, storage.ErrAlreadyExist) {
			log.Error(err.Error(), slog.String("email", input.Email))

			return nil, custom_errors.RegistrationDuplicateError
		}
		log.Error(err.Error())
	}

	token, refresh, err := jwt_helper.GenerateTokens(dbUser.Id)
	if err != nil {
		log.Error(err.Error(), slog.String("email", input.Email))

		return nil, custom_errors.AuthenticationError
	}

	return &dto.TokenResponse{Token: token, RefreshToken: refresh}, nil
}

func (s *JwtAuthService) GetUser(userId uuid.UUID) (*dto.GetUserResponse, error) {
	const op = "JwtAuthService.GetUser"
	log := logging.CreateLoggerWithOp(op)

	user, err := s.userStore.GetUser(userId.String())
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("Пользователь не найден", err)
		}

		return nil, errors.Join(errors.New("Пользователь не найден"), err)
	}

	response := dto.GetUserResponse{
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role}
	return &response, nil
}

func (s *JwtAuthService) Login(input LoginInput) (*dto.TokenResponse, error) {
	const op = "JwtAuthService.Login"
	log := logging.CreateLoggerWithOp(op)

	hasher := hash.BcryptHash{}
	passwordHash := hash.HashPassword(hasher, input.Password)

	user, err := s.userStore.CheckCredentials(input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrNoUserWithCred) {
			log.Warn("Пользователя с такой почтой нет", slog.String("email", input.Email))

			return nil, custom_errors.NoUserError
		} else if errors.Is(err, storage.ErrPassIncorrect) {
			log.Warn("Неверный пароль", slog.String("email", input.Email))

			return nil, custom_errors.InvalidCredError
		}

		return nil, custom_errors.InternalError
	}

	token, refresh, err := jwt_helper.GenerateTokens(user.Id)
	if err != nil {
		log.Error(
			err.Error(),
			slog.String("additional", "ошибка при генерации токенов"),
			slog.String("email", input.Email))

		return nil, custom_errors.InternalError
	}

	return &dto.TokenResponse{
		Token:        token,
		RefreshToken: refresh,
	}, nil

}

func (s *JwtAuthService) Reauth(refresh string) (string, error) {
	const op = "JwtAuthService.Reauth"
	log := logging.CreateLoggerWithOp(op)

	newToken, _, err := jwt_helper.Reauth(refresh)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка при регенерации токена: %s", err.Error()))

		return "", err
	}

	return newToken, nil
}
