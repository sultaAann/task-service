package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Description string    `json:"description"`
}

type CreatedResponse struct {
	ID          int       `json:"id"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
}

type Status int

const (
	Pending Status = iota
	In_progress
	Done
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case In_progress:
		return "In_progress"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
func (s *Status) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch j {
	case "Pending":
		*s = Pending
	case "In_progress":
		*s = In_progress
	case "Done":
		*s = Done
	default:
		return fmt.Errorf("invalid TaskState value: %s", j)
	}
	return nil
}
