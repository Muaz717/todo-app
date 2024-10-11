package item

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Muaz717/todo-app/internal/app/http-server/middleware/identification"
	"github.com/Muaz717/todo-app/internal/domain/models"
	resp "github.com/Muaz717/todo-app/internal/lib/api/response"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.2 --name=Item
type Item interface {
	Create(
		ctx context.Context,
		userId int64,
		title string,
		description string,
	) (int64, error)
	AllItems(ctx context.Context, userId int64) ([]models.Item, error)
}

type ItemHandler struct {
	ctx  context.Context
	log  *slog.Logger
	item Item
}

func New(
	ctx context.Context,
	log *slog.Logger,
	item Item,
) *ItemHandler {
	return &ItemHandler{
		ctx:  ctx,
		log:  log,
		item: item,
	}
}

type Request struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.item.Create"

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

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(validateErr))

		return
	}

	userId, err := identification.GetUserId(r)
	if err != nil {
		log.Error("failed to get user id", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to get user id"))

		return
	}

	itemId, err := h.item.Create(h.ctx, userId, req.Title, req.Description)
	if err != nil {
		log.Error("failed to create item", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to create item"))

		return
	}

	log.Info("item created", slog.Int64("item_id", itemId), slog.Int64("user_id", userId))

	render.JSON(w, r, resp.OK("Item successfully created"))
}

func (h *ItemHandler) AllItems(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.item.AllItems"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userId, err := identification.GetUserId(r)
	if err != nil {
		log.Error("failed to get user id", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("failed to get user id"))

		return
	}

	items, err := h.item.AllItems(h.ctx, userId)
	if err != nil {
		log.Error("failed to get items", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to get items")) //TODO ...

		return
	}

	log.Info("All lists showed")

	render.JSON(w, r, items)
}
