package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"task-sevice/internal/custom_errors"
	"task-sevice/internal/repository"
	"task-sevice/internal/service"
)

var (
	TaskURL    = regexp.MustCompile(`^/task/*$`)
	TaskWithID = regexp.MustCompile(`^/task/([1-9][0-9]*)$`)
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	s service.Service
}

func NewHandler(s service.Service) Handler {
	return &handler{s: s}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodDelete && TaskWithID.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
	case r.Method == http.MethodGet && TaskWithID.MatchString(r.URL.Path):
		h.GetById(w, r)
		return
	case r.Method == http.MethodPost && TaskURL.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	// case r.Method == http.MethodPut && TaskURL.MatchString(r.URL.Path):
	// 	h.Update(w, r)
	// 	return
	case r.Method == http.MethodGet && TaskURL.MatchString(r.URL.Path):
		h.GetAll(w, r)
		return
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks := h.s.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	matches := TaskWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	task, err := h.s.GetById(id)

	if err != nil {
		if errors.As(err, &custom_errors.NotFoundError{}) {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
}
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var description map[string]interface{}
	var task repository.Task
	if err := json.NewDecoder(r.Body).Decode(&description); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	task.Description = description["description"].(string)
	id := h.s.Create(task)
	t, err := h.s.GetById(id)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	var result repository.CreatedResponse = repository.CreatedResponse{
		ID:          id,
		CreatedAt:   t.CreatedAt,
		Status:      t.Status,
		Description: t.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
}

// func (h *handler) Update(w http.ResponseWriter, r *http.Request)  {}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := TaskWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	err = h.s.DeleteById(id)

	if err != nil {
		if errors.As(err, &custom_errors.NotFoundError{}) {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": "500 Internal Server Error"})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "404 Not Found: Check Path"})
}
