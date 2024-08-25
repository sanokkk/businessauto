package storage

import "errors"

var (
	ErrPrepareQuery = errors.New("Ошибка в формировании SQL-запроса")
	ErrInsert       = errors.New("Ошибка во вставке в БД")
	ErrAlreadyExist = errors.New("Такая запись уже есть в БД")
	ErrNotFound     = errors.New("Запись не найдена")
	ErrSearch       = errors.New("Ошибка поиска сущности")

	ErrNoUserWithCred = errors.New("Нет пользователя с таким логином")
	ErrPassIncorrect  = errors.New("Неверный пароль")
)
