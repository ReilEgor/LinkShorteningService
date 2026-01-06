package usecase

import (
	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
)

type TaskUsecase struct {
	repo domain.TaskRepository
}

func NewTaskUsecase(r domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: r}
}

func (uc *TaskUsecase) AddTask(title string) (*domain.Task, error) {
	if title == "" {
		return nil, domain.ErrEmptyTitle
	}

	newTask := &domain.Task{Title: title, Done: false}
	err := uc.repo.Create(newTask)
	return newTask, err
}

func (uc *TaskUsecase) GetTasks() ([]domain.Task, error) {
	return uc.repo.GetAll()
}
