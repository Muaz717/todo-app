package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	resp "github.com/Muaz717/todo-app/internal/lib/api/response"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.2 --name=Auth
type Auth interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
}

type AuthHandler struct {
	ctx  context.Context
	log  *slog.Logger
	auth Auth
}

func New(ctx context.Context, log *slog.Logger, auth Auth) *AuthHandler {
	return &AuthHandler{
		ctx:  ctx,
		log:  log,
		auth: auth,
	}
}

func (h *AuthHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.RegisterNewUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("empty request"))

		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}

	log.Info("request body decoded", slog.Any("req", req))

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(validateErr))

		return
	}

	userId, err := h.auth.RegisterNewUser(h.ctx, req.Email, req.Password)
	if err != nil {
		log.Error("failed to register new user", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to register new user"))

		return
	}

	log.Info("User successfully registered", slog.Int("userId", int(userId)))

	render.JSON(w, r, resp.OK("You successfully registered"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("empty request"))

		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}

	log.Info("request body decoded", slog.Any("req", req))

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(validateErr))

		return
	}

	token, err := h.auth.Login(h.ctx, req.Email, req.Password)
	if err != nil {
		log.Error("invalid email or password", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("invalid email or password"))

		return
	}

	log.Info("user got token", slog.String("token", token))

	render.JSON(w, r, responseOK(token))
}

func responseOK(token string) Response {
	return Response{
		Response: resp.OK("You got token"),
		Token:    token,
	}
}
