package storage

import (
	"autoshop/internal/config"
	"autoshop/internal/domain/models"
	"autoshop/pkg/hash"
	"autoshop/pkg/logging"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

type UsersStorage struct {
	db *gorm.DB
}

func NewUsersStorage(cfg *config.DbConfig) *UsersStorage {
	db, err := gorm.Open(postgres.Open(cfg.DbConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	return &UsersStorage{db: db}
}

/*SaveUser(user *models.User) error
GetUser(id string) (*models.User, error)
CheckCredentials(email string, password string) (*models.User, error)*/

func (r *UsersStorage) SaveUser(user *models.User) error {
	const op = "UsersStorage.SaveUser"
	log := logging.CreateLoggerWithOp(op)

	res := r.db.Table("users").Create(user)
	if res.Error != nil {

		var pgerror *pgconn.PgError
		if errors.As(res.Error, &pgerror) {
			if pgerror.Code == "23505" {
				log.Warn("Пользователь с данной почтой уже существует", slog.String("email", user.Email))

				return ErrAlreadyExist
			}
		}

		log.Warn("Ошибка при вставке пользователя в БД", slog.Any("user", user), slog.String("error", res.Error.Error()))

		return res.Error
	}

	return nil
}

func (r *UsersStorage) GetUser(id uuid.UUID) (*models.User, error) {
	const op = "UsersStorage.GetUser"
	log := logging.CreateLoggerWithOp(op)

	var user models.User
	res := r.db.Table("users").First(&models.User{Id: id}).Find(&user)
	if res.Error != nil {
		log.Warn("Ошибка при поиске пользователя в БД", slog.Any("userId", id), slog.String("error", res.Error.Error()))

		return nil, res.Error
	}

	return &user, nil
}

func (r *UsersStorage) CheckCredentials(email string, password string) (*models.User, error) {
	const op = "UsersStorage.CheckCredentials"
	log := logging.CreateLoggerWithOp(op)

	var user models.User
	res := r.db.Table("users").Where("email = ?", email).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Error("Нет пользователя с такой почтой", slog.String("email", email))

			return nil, ErrNoUserWithCred
		}

		log.Warn("Ошибка при поиске пользователя в БД", slog.Any("email", email), slog.String("error", res.Error.Error()))

		return nil, res.Error
	}

	hasher := hash.BcryptHash{}
	if !hash.ComparePasswordAndHash(hasher, password, user.PasswordHash) {
		log.Error("Не совпадают пароли", slog.String("email", email))

		return nil, ErrPassIncorrect
	}

	return &user, nil
}
