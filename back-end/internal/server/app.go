package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/config"
	"github.com/ReilEgor/CleanArchitectureGolang/internal/repository/postgres"

	api "github.com/ReilEgor/CleanArchitectureGolang/internal/delivery/http"
	"github.com/ReilEgor/CleanArchitectureGolang/internal/usecase"
)

func initDB(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := initDB(cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("postgres connection failed: %w", err)
	}
	defer db.Close()

	repo := postgres.NewTaskRepo(db)

	uc := usecase.NewTaskUsecase(repo)

	handler := api.NewGinServer(uc)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: handler.GetRouter(),
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("Server is starting on", "port", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("Start shutdown", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	slog.Info("Server stopped correctly")
	return nil
}
