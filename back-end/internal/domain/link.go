package domain

import (
	"context"
	"errors"
)

type Link struct {
	ID       string `json:"id"`
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
}

var (
	ErrLinkNotFound     = errors.New("link not found")
	ErrInvalidURL       = errors.New("url is unreachable")
	ErrInvalidShortCode = errors.New("invalid short code")
)

type LinkRepository interface {
	Create(ctx context.Context, link *Link) (uint64, error)
	Get(ctx context.Context, id uint64) (Link, error)
}
type LinkUsecase interface {
	AddLink(ctx context.Context, url string) (*Link, error)
	GetLink(ctx context.Context, url string) (Link, error)
}
