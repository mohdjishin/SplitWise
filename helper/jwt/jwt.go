package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohdjishin/SplitWise/config"
)

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(config.GetConfig().JwtString))
}
