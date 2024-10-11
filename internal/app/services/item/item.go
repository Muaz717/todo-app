package itemsrv

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Muaz717/todo-app/internal/domain/models"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
)

type Item struct {
	log *slog.Logger
	ItemSaver
	ItemProvider
}

type ItemSaver interface {
	SaveItem(
		ctx context.Context,
		userId int64,
		title string,
		description string,
	) (int64, error)
}

type ItemProvider interface {
	AllItems(ctx context.Context, userId int64) ([]models.Item, error)
}

func New(
	log *slog.Logger,
	itemSaver ItemSaver,
	itemProvider ItemProvider,
) *Item {
	return &Item{
		log:          log,
		ItemSaver:    itemSaver,
		ItemProvider: itemProvider,
	}
}

func (i *Item) Create(
	ctx context.Context,
	userId int64,
	title string,
	description string,
) (int64, error) {
	const op = "services.item.Create"

	log := i.log.With(
		slog.String("op", op),
	)

	log.Info("Creating item")

	itemId, err := i.ItemSaver.SaveItem(ctx, userId, title, description)
	if err != nil {
		log.Error("failed to save item")

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("item saved", slog.Int64("id", itemId))

	return itemId, nil
}

func (i *Item) AllItems(ctx context.Context, userId int64) ([]models.Item, error) {
	const op = "services.item.AllItems"

	log := i.log.With(
		slog.String("op", op),
	)

	log.Info("Getting items")

	items, err := i.ItemProvider.AllItems(ctx, userId)
	if err != nil {
		log.Error("failed to got items", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Got items")

	return items, nil
}
