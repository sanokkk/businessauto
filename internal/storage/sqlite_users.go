package storage

import (
	"autoshop/internal/domain/models"
	"autoshop/pkg/hash"
	"autoshop/pkg/logging"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"log/slog"
)

type SqliteUserStorage struct {
	db *sql.DB
}

func NewSqliteUserStorage(db *sql.DB) *SqliteUserStorage {
	return &SqliteUserStorage{db: db}
}

func (r *SqliteUserStorage) SaveUser(user *models.User) error {
	const op = "SqliteUserStorage.SaveUser"
	log := logging.CreateLoggerWithOp(op)

	insertUser := `
INSERT INTO users (id, email, fullName, passwordHash) 
VALUES (
        $1, $2, $3, $4
)
`
	stmt, err := r.db.Prepare(insertUser)
	if err != nil {
		log.Error(err.Error())

		return ErrPrepareQuery
	}

	if _, err := stmt.Exec(user.Id, user.Email, user.FullName, user.PasswordHash); err != nil {
		log.Error(err.Error())

		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr, sqlite3.ErrConstraintUnique) {
				return ErrAlreadyExist
			}
		}
		return ErrInsert
	}

	return nil
}

func (r *SqliteUserStorage) GetUser(id string) (*models.User, error) {
	const op = "SqliteUserStorage.GetUser"
	log := logging.CreateLoggerWithOp(op)

	getUser := `
SELECT id, email, fullName, role FROM users
WHERE id=$1
LIMIT 1
`
	row := r.db.QueryRow(getUser, id)

	var user models.User
	var userId string

	if err := row.Scan(&userId, &user.Email, &user.FullName, &user.Role); err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if errors.Is(sqliteError, sqlite3.ErrNotFound) {
				log.Error(fmt.Sprintf("Не найдена запись с Id %s", id))

				return nil, ErrNotFound
			}
		}

		log.Error(err.Error())
		return nil, errors.Join(ErrSearch, err)
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	user.Id = uid
	return &user, nil
}

func (r *SqliteUserStorage) CheckCredentials(email string, password string) (*models.User, error) {
	const op = "SqliteUserStorage.CheckCredentials"
	log := logging.CreateLoggerWithOp(op)

	getUser := `
SELECT id, passwordHash FROM users
WHERE email=$1
LIMIT 1
`
	row := r.db.QueryRow(getUser, email)

	var user models.User
	var userId string

	if err := row.Scan(&userId, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Нет пользователя", slog.String("email", email))

			return nil, ErrNoUserWithCred
		}

		log.Error(err.Error())
		return nil, errors.Join(ErrSearch, err)
	}

	hasher := hash.BcryptHash{}
	if !hash.ComparePasswordAndHash(hasher, password, user.PasswordHash) {
		log.Error("Не совпадают пароли", slog.String("email", email))

		return nil, ErrPassIncorrect
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	user.Id = uid
	return &user, nil
}
