package model

import "errors"

var (
	ErrEmptyText = errors.New("текст задачи не может быть пустым")
	ErrNotFound  = errors.New("задача не найдена")
)
