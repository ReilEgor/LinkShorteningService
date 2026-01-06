package repository

import (
	"sync"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
)

type MemoryTaskRepo struct {
	mu    sync.Mutex
	tasks []domain.Task
}

func NewMemoryTaskRepo() *MemoryTaskRepo {
	return &MemoryTaskRepo{tasks: []domain.Task{}}
}

func (r *MemoryTaskRepo) Create(t *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t.ID = len(r.tasks) + 1
	r.tasks = append(r.tasks, *t)
	return nil
}

func (r *MemoryTaskRepo) GetAll() ([]domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tasks, nil
}
