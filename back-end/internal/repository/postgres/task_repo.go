package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ReilEgor/LinkShorteningService/internal/domain"
)

type LinkRepo struct {
	db *sql.DB
}

func NewLinkRepo(db *sql.DB) *LinkRepo {
	return &LinkRepo{db: db}
}

const createLinkQuery = `INSERT INTO links (long_url) VALUES ($1) RETURNING id`

func (r *LinkRepo) Create(ctx context.Context, link *domain.Link) (uint64, error) {
	var id uint64
	err := r.db.QueryRowContext(ctx, createLinkQuery, link.LongURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo create link: %w", err)
	}

	return id, nil
}

const getLinkQuery = `SELECT id, long_url FROM links WHERE id = $1`

func (r *LinkRepo) Get(ctx context.Context, id uint64) (domain.Link, error) {
	var l domain.Link
	err := r.db.QueryRowContext(ctx, getLinkQuery, id).Scan(&l.ID, &l.LongURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, nil
		}
		return domain.Link{}, fmt.Errorf("repo get link: %w", err)
	}

	return l, nil
}
