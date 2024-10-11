package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/auth"
	"github.com/Muaz717/todo-app/internal/app/http-server/handlers/item"
	"github.com/Muaz717/todo-app/internal/app/http-server/middleware/identification"
	mwLogger "github.com/Muaz717/todo-app/internal/app/http-server/middleware/logger"

	"github.com/Muaz717/todo-app/internal/config"
	"github.com/Muaz717/todo-app/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	HTTPServer *http.Server
	router     chi.Router
	ctx        context.Context
	log        *slog.Logger
	cfg        config.Config
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg config.Config,
	authSrv auth.Auth,
	itemSrv item.Item,
) *App {

	authHandler := auth.New(ctx, log, authSrv)
	itemHandler := item.New(ctx, log, itemSrv)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func(auth chi.Router) {
		auth.Post("/sign-up", authHandler.RegisterNewUser)
		auth.Post("/sign-in", authHandler.Login)
	})

	router.Route("/api", func(api chi.Router) {
		api.Use(identification.New(log))

		api.Route("/items", func(items chi.Router) {
			items.Post("/", itemHandler.Create)
			items.Get("/", itemHandler.AllItems)
		})
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		HTTPServer: srv,
		router:     router,
		ctx:        ctx,
		log:        log,
		cfg:        cfg,
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("addr", a.cfg.Address),
	)

	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Error("failed to run http server", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("HTTP server is running", slog.String("addr", a.HTTPServer.Addr))

	return nil
}

func (a *App) Stop() error {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping Http server", slog.String("addr", a.HTTPServer.Addr))

	if err := a.HTTPServer.Shutdown(a.ctx); err != nil {
		a.log.Error("failed to stop server", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
