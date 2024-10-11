package app

import (
	"context"
	"log/slog"

	authService "github.com/Muaz717/todo-app/internal/app/services/auth"
	itemsrv "github.com/Muaz717/todo-app/internal/app/services/item"
	"github.com/Muaz717/todo-app/internal/app/storage/postgres"
	"github.com/Muaz717/todo-app/internal/config"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
	httpapp "github.com/Muaz717/todo-app/internal/pkg/app/http"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		panic(err)
	}

	authSrv := authService.New(log, storage, storage, cfg.TokenTTL)
	itemSrv := itemsrv.New(log, storage, storage)

	httpApp := httpapp.New(ctx, log, *cfg, authSrv, itemSrv)

	return &App{
		HTTPSrv: httpApp,
	}
}
