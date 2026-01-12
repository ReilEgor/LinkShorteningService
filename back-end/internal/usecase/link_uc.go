package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ReilEgor/LinkShorteningService/internal/domain"
	"github.com/ReilEgor/LinkShorteningService/pkg/shortener"
)

type LinkUsecase struct {
	repo   domain.LinkRepository
	logger *slog.Logger
}

func NewLinkUsecase(r domain.LinkRepository, l *slog.Logger) *LinkUsecase {
	return &LinkUsecase{
		repo:   r,
		logger: l,
	}
}

func (uc *LinkUsecase) AddLink(ctx context.Context, url string) (*domain.Link, error) {
	if !uc.isReachableURL(ctx, url) {
		return nil, errors.New("URL is not reachable or timed out")
	}
	link := &domain.Link{
		LongURL: url,
	}
	id, err := uc.repo.Create(ctx, link)
	if err != nil {
		uc.logger.Error("failed to create link in repo", slog.Any("err", err))
		return nil, fmt.Errorf("repository error: %w", err)
	}
	link.ID = strconv.FormatUint(id, 10)

	link.ShortURL = shortener.Encode(id)

	return link, nil
}

func (uc *LinkUsecase) GetLink(ctx context.Context, shortCode string) (domain.Link, error) {
	id, err := shortener.Decode(shortCode)
	if err != nil {
		uc.logger.Warn("failed to decode short code", slog.String("code", shortCode))
		return domain.Link{}, fmt.Errorf("invalid short code: %w", err)
	}
	link, err := uc.repo.Get(ctx, id)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get link: %w", err)
	}

	return link, nil
}

func (uc *LinkUsecase) isReachableURL(ctx context.Context, url string) bool {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(checkCtx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		uc.logger.Debug("url unreachable", slog.String("url", url), slog.Any("err", err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
