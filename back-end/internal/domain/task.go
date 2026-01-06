package domain

import "errors"

var ErrEmptyTitle = errors.New("title cannot be empty")

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskRepository interface {
	Create(task *Task) error
	GetAll() ([]Task, error)
}
type TaskUsecase interface {
	AddTask(title string) (*Task, error)
	GetTasks() ([]Task, error)
}
