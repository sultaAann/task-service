package repository

import "task-sevice/internal/custom_errors"

type Repository interface {
	GetAll() []Task
	GetById(id int) (*Task, error)
	Create(t Task) int
	Update(id int, t Task) (*Task, error)
	DeleteById(id int) error
}

type repository struct {
	db *db
}

func NewRepository(db *db) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() []Task {
	return r.db.GetAll()
}

func (r *repository) GetById(id int) (*Task, error) {
	t, exists := r.db.GetTask(id)
	result := &Task{
		ID:          t.ID,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		StartedAt:   t.StartedAt,
		CompletedAt: t.CompletedAt,
	}
	if !exists {
		return nil, custom_errors.NewNotFoundError("DB:", id, "Task Not Found")
	}
	return result, nil
}

func (r *repository) Create(t Task) int {
	return r.db.AddTask(t)
}

func (r *repository) Update(id int, t Task) (*Task, error) {
	exists := r.db.UpdateTask(id, t)
	if !exists {
		return nil, custom_errors.NewNotFoundError("DB:", t.ID, "Task Not Found")
	}
	return r.GetById(t.ID)
}

func (r *repository) DeleteById(id int) error {
	exists := r.db.DeleteTask(id)
	if !exists {
		return custom_errors.NewNotFoundError("DB:", id, "Task Not Found")
	}
	return nil
}
