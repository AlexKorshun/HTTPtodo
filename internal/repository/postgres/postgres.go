package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(databaseURL string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{pool: pool}, nil
}

func (s *PostgresStorage) GetAll() ([]model.Task, error) {
	ctx := context.Background()
	tasks := []model.Task{}
	rows, err := s.pool.Query(ctx, "SELECT id, text, done FROM tasks")
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

func (s *PostgresStorage) Create(text string) (model.Task, error) {
	ctx := context.Background()
	t := model.Task{}
	row := s.pool.QueryRow(ctx, "INSERT INTO tasks (text) VALUES ($1) RETURNING id, text, done", text)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if err != nil {
		return model.Task{}, fmt.Errorf("Create: %w", err)
	}
	return t, nil
}

func (s *PostgresStorage) ToggleDone(id int) (model.Task, error) {
	ctx := context.Background()
	t := model.Task{}
	row := s.pool.QueryRow(ctx, "UPDATE tasks SET done = NOT done WHERE id = $1 RETURNING id, text, done", id)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Task{}, model.ErrNotFound
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: %w", err)
	}
	return t, nil
}

func (s *PostgresStorage) Delete(id int) error {
	ctx := context.Background()
	tag, err := s.pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil

}
