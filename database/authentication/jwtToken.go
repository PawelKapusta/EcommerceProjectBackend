package authentication

import (
	"backend/database/models"
	"backend/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func GetCustomClaimsConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:      &JwtCustomClaims{},
		SigningKey:  []byte(utils.GetValueFromEnv("JWT_SECRET", "secret")),
		TokenLookup: "header:Authorization",
	}
}

type JwtCustomClaims struct {
	models.User
	jwt.StandardClaims
}

func CreateToken(userData models.User) (string, error) {
	claims := &JwtCustomClaims{
		models.User{
			Name:  userData.Name,
			Email: userData.Email,
			Admin: userData.Admin,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(utils.GetValueFromEnv("JWT_SECRET", "secret")))
	if err != nil {
		return "", err
	}
	return t, nil
}
