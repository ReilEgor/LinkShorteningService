package postgres

import (
	"database/sql"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(t *domain.Task) error {
	query := "INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRow(query, t.Title, t.Done).Scan(&t.ID)
}

func (r *TaskRepo) GetAll() ([]domain.Task, error) {
	rows, err := r.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
