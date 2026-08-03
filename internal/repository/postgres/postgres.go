package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
	"github.com/AlexKorshun/HTTPtodo/internal/service"
)

const (
	queryGetAll     = `SELECT id, text, done FROM tasks`
	queryCreate     = `INSERT INTO tasks (text) VALUES ($1) RETURNING id, text, done`
	queryToggleDone = `UPDATE tasks SET done = NOT done WHERE id = $1 RETURNING id, text, done`
	queryDelete     = `DELETE FROM tasks WHERE id = $1`
)

type Storage struct {
	pool *pgxpool.Pool
}

var _ service.Storage = (*Storage)(nil)

func New(databaseURL string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Task, error) {
	tasks := []model.Task{}
	rows, err := s.pool.Query(ctx, queryGetAll)
	if err != nil {
		return tasks, fmt.Errorf("GetAll: чтение из базы даных: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		t := model.Task{}
		err := rows.Scan(&t.ID, &t.Text, &t.Done)
		if err != nil {
			return []model.Task{}, fmt.Errorf("GetAll: ошибка внутри цикла чтения на элементе ID = %d: %w", t.ID, err)
		}

		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return []model.Task{}, fmt.Errorf("GetAll: ошибка после цикла: %w", err)
	}
	return tasks, nil
}

func (s *Storage) Create(ctx context.Context, text string) (model.Task, error) {
	t := model.Task{}
	row := s.pool.QueryRow(ctx, queryCreate, text)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if err != nil {
		return model.Task{}, fmt.Errorf("Create: %w", err)
	}
	return t, nil
}

func (s *Storage) ToggleDone(ctx context.Context, id int) (model.Task, error) {
	t := model.Task{}
	row := s.pool.QueryRow(ctx, queryToggleDone, id)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Task{}, model.ErrNotFound
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: %w", err)
	}
	return t, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	tag, err := s.pool.Exec(ctx, queryDelete, id)

	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil

}
