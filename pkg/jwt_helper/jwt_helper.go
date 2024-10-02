package jwt_helper

import (
	"autoshop/internal/config"
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type SignedDetails struct {
	//Email    string
	//FullName string
	UserId string
	jwt.StandardClaims
	Role string
}

func GenerateTokens(uid uuid.UUID, role string) (signedToken string, refreshToken string, err error) {
	const op = "JwtHelper.GenerateTokens"
	logger := logging.CreateLoggerWithOp(op)

	jwtConfig := config.MustLoadConfig().JwtConfig

	signedToken, err = generateToken(uid.String(), role, jwtConfig.ExpireAfter, jwtConfig.Secret)
	if err != nil {
		logger.Error(err.Error())

		return "", "", err
	}

	refreshToken, err = generateToken(uid.String(), role, jwtConfig.RefreshExpireAfter, jwtConfig.Secret)
	if err != nil {
		logger.Error(err.Error())

		return "", "", err
	}

	return
}

func ValidateToken(token string) (claims *SignedDetails, isValid bool, errMsg error) {
	jwtConfig := config.MustLoadConfig().JwtConfig

	jwtToken, err := jwt.ParseWithClaims(
		token,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

	if err != nil {
		errMsg = err
		return
	}

	claims, ok := jwtToken.Claims.(*SignedDetails)
	if !ok {
		errMsg = fmt.Errorf("Ошибка парсинга токена: %s", err.Error())
		return
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		errMsg = custom_errors.TokenExpiredError
		return
	}

	return claims, true, nil
}

func Reauth(refresh string) (jwtToken string, refreshToken string, reauthError error) {
	const op = "JwtHelper.Reauth"
	log := logging.CreateLoggerWithOp(op)

	claims := SignedDetails{}
	jwtConfig := config.MustLoadConfig().JwtConfig

	_, isRefreshValid, errMsg := ValidateToken(refresh)
	if !isRefreshValid {
		msg := fmt.Sprintf("Ошибка при рефреше: %s", errMsg)
		log.Error(msg)

		return "", "", errors.New(msg)
	}

	_, err := jwt.ParseWithClaims(
		refresh,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

	if err != nil {
		log.Error(err.Error())

		return "", "", err
	}

	uid := claims.UserId
	role := claims.Role
	jwtToken, err = generateToken(uid, role, jwtConfig.ExpireAfter, jwtConfig.Secret)
	if err != nil {
		log.Error(err.Error())

		return "", "", err
	}

	return jwtToken, refresh, nil
}

// todo add roles
func generateToken(userId string, role string, expireAfter time.Duration, secret string) (string, error) {
	claims := &SignedDetails{
		UserId: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(expireAfter).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
