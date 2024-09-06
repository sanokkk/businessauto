package custom_errors

import "errors"

var (
	ValidationError            = errors.New("Ошибка валидации. Введенные данные некорректны")
	RegistrationError          = errors.New("Ошибка регистрации. Попробуйте еще раз")
	RegistrationDuplicateError = errors.New("Пользователь с таким Email уже существует")
	AuthenticationError        = errors.New("Ошибка аутентификации")
	TokenExpiredError          = errors.New("Токен протух")
	NoUserError                = errors.New("Нет пользователя с такими данными")
	InvalidCredError           = errors.New("Неверный логин/пароль")
	InternalError              = errors.New("Неизвестная ошибка")
	NoSuchFilterError          = errors.New("Нет такого фильтра")
	ConvertationError          = errors.New("Ошибка конвертации")
	NoFileError                = errors.New("Данного файла не существует")
)
