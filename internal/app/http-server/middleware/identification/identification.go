package identification

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	resp "github.com/Muaz717/todo-app/internal/lib/api/response"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeader = "Authorization"
)

type Uid string

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		const op = "middleware.Identification.New"

		log := log.With(
			slog.String("component", "user identification"),
			slog.String("op", op),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(authorizationHeader)

			token, err := validateToken(log, header)
			if err != nil {
				log.Error("failed to validate token", sl.Err(err))

				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error(err.Error()))

				return
			}

			tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("MY_SECRET")), nil
			})
			if err != nil {
				log.Error("failed to parse token", sl.Err(err))

				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("failed to parse token"))

				return
			}

			claims := tokenParsed.Claims.(jwt.MapClaims)

			userId := int64(claims["uid"].(float64))

			log.Info("token successfully parsed")

			uidStr := Uid("user_id")

			ctx := r.Context()
			withValue := context.WithValue(ctx, uidStr, userId)
			r = r.WithContext(withValue)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func GetUserId(r *http.Request) (int64, error) {
	userId := r.Context().Value(Uid("user_id"))
	if userId == "" {
		return 0, errors.New("user id is empty")
	}

	idInt, ok := userId.(int64)
	if !ok {
		return 0, errors.New("user id is invalid type")
	}

	return idInt, nil
}

func validateToken(log *slog.Logger, header string) (string, error) {
	if header == "" {
		log.Error("empty authorization header")

		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Error("invalid auth token")

		return "", errors.New("invalid auth token")
	}

	if len(headerParts[1]) == 0 {
		log.Error("token is empty")

		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
