package service

import (
	"task-sevice/internal/repository"
	"time"
)

type Service interface {
	GetAll() []repository.Task
	GetById(id int) (*repository.Task, error)
	Create(t repository.Task) int
	// Update()
	DeleteById(id int) error
}

type service struct {
	r repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{r}
}

func (s service) GetAll() []repository.Task {
	return s.r.GetAll()
}

func (s service) GetById(id int) (*repository.Task, error) {
	return s.r.GetById(id)
}

func (s service) Create(t repository.Task) int {
	t.CreatedAt = time.Now()
	t.Status = repository.Pending

	id := s.r.Create(t)

	go s.work(id)

	return id
}

func (s service) DeleteById(id int) error {
	return s.r.DeleteById(id)
}

func (s service) work(id int) {
	time.Sleep(1 * time.Second)

	t, err := s.r.GetById(id)
	if err != nil {
		return
	}
	t.StartedAt = time.Now()
	t.Status = repository.In_progress
	s.r.Update(id, *t)

	time.Sleep(3 * time.Minute)

	t.CompletedAt = time.Now()
	t.Status = repository.Done
	s.r.Update(id, *t)
}
