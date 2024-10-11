package jwt

import (
	"fmt"
	"time"

	"github.com/Muaz717/todo-app/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	const op = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
