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

	"github.com/ReilEgor/LinkShorteningService/internal/config"
	apigRPC "github.com/ReilEgor/LinkShorteningService/internal/delivery/gRPC"
	"github.com/ReilEgor/LinkShorteningService/internal/repository/postgres"
	"github.com/ReilEgor/LinkShorteningService/pkg/logger"

	apiHTTP "github.com/ReilEgor/LinkShorteningService/internal/delivery/http"
	"github.com/ReilEgor/LinkShorteningService/internal/usecase"
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

	h := slog.NewJSONHandler(os.Stdout, nil)
	myLogger := slog.New(&logger.ContextHandler{Handler: h})
	slog.SetDefault(myLogger)

	db, err := initDB(cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			myLogger.Error("database close failed:", err)
		}
	}(db)

	repo := postgres.NewLinkRepo(db)

	uc := usecase.NewLinkUsecase(repo, myLogger)

	handler := apiHTTP.NewGinServer(uc, myLogger)

	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handler.GetRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := apigRPC.RunGRPCServer(ctx, "50051", uc, myLogger); err != nil {
			serverErrors <- fmt.Errorf("grpc: %w", err)
		}
	}()

	go func() {
		slog.Info("server is starting",
			slog.String("addr", httpServer.Addr),
		)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("listen and serve: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("critical error: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	slog.Info("stopping servers...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		myLogger.Error("http shutdown failed", "error", err)
		return fmt.Errorf("could not stop http server: %w", err)
	}

	slog.Info("Server stopped correctly")
	return nil
}
