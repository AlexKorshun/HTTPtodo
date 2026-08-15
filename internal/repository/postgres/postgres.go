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

func (s *Storage) GetAll(ctx context.Context, userID int64) ([]model.Task, error) {
	const queryGetAll = `SELECT id, text, done FROM tasks WHERE user_id = $1`
	tasks := []model.Task{}
	rows, err := s.pool.Query(ctx, queryGetAll, userID)
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

func (s *Storage) Create(ctx context.Context, userID int64, text string) (model.Task, error) {
	const queryCreate = `INSERT INTO tasks (text, user_id) VALUES ($1, $2) RETURNING id, text, done`
	t := model.Task{}
	row := s.pool.QueryRow(ctx, queryCreate, text, userID)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if err != nil {
		return model.Task{}, fmt.Errorf("Create: %w", err)
	}
	return t, nil
}

func (s *Storage) ToggleDone(ctx context.Context, userID int64, id int) (model.Task, error) {
	const queryToggleDone = `UPDATE tasks SET done = NOT done WHERE id = $1 AND user_id = $2 RETURNING id, text, done`
	t := model.Task{}
	row := s.pool.QueryRow(ctx, queryToggleDone, id, userID)
	err := row.Scan(&t.ID, &t.Text, &t.Done)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Task{}, model.ErrNotFound
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("ToggleDone: %w", err)
	}
	return t, nil
}

func (s *Storage) Delete(ctx context.Context, userID int64, id int) error {
	const queryDelete = `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	tag, err := s.pool.Exec(ctx, queryDelete, id, userID)

	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil

}
