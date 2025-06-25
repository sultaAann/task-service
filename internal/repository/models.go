package repository

import "time"

type Task struct {
	ID          int
	Status      Status // pending, in_progress, done
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
	Description string
}

type Status int

const (
	Pending Status = iota
	In_progress
	Done
)

type CreatedResponse struct {
	ID          int
	Status      Status
	CreatedAt   time.Time
	Description string
}
