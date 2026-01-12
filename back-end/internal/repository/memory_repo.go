package repository

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/ReilEgor/LinkShorteningService/internal/domain"
)

type MemoryLinkRepo struct {
	mu     sync.RWMutex
	links  map[uint64]domain.Link
	nextID uint64
}

func NewMemoryLinkRepo() *MemoryLinkRepo {
	return &MemoryLinkRepo{
		links:  make(map[uint64]domain.Link),
		nextID: 1,
	}
}

func (r *MemoryLinkRepo) Create(ctx context.Context, t *domain.Link) (uint64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	t.ID = strconv.FormatUint(id, 10)

	r.links[id] = *t
	r.nextID++

	return id, nil
}

func (r *MemoryLinkRepo) Get(ctx context.Context, id uint64) (domain.Link, error) {
	select {
	case <-ctx.Done():
		return domain.Link{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.links[id]
	if !ok {
		return domain.Link{}, errors.New("link not found in database")
	}

	return link, nil
}
