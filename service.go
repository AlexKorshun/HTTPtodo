package main

type Storage interface {
	GetAll() ([]Task, error)
	Create(text string) (Task, error)
	ToggleDone(id int) (Task, error)
	Delete(id int) error
}

type TaskService struct {
	storage Storage
}

func (a *TaskService) List() ([]Task, error) {
	return a.storage.GetAll()
}

func (a *TaskService) Add(text string) (Task, error) {
	if text == "" {
		return Task{}, ErrEmptyText
	}
	return a.storage.Create(text)
}

func (a *TaskService) Change(id int) (Task, error) {
	return a.storage.ToggleDone(id)
}

func (a *TaskService) Delete(id int) error {
	return a.storage.Delete(id)

}
