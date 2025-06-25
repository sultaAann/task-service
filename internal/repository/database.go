package repository

import "sync"

var (
	instance *db
	once     sync.Once
)

type db struct {
	tasks  map[int]Task
	nextID int
	mu     sync.RWMutex
}

func NewInstanceDB() *db {
	once.Do(func() {
		instance = &db{
			tasks:  make(map[int]Task),
			nextID: 1,
		}
	})
	return instance
}

func (d *db) GetAll() []Task {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := []Task{}
	for _, value := range d.tasks {
		result = append(result, value)
	}
	return result
}

func (d *db) AddTask(task Task) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	task.ID = d.nextID
	d.tasks[d.nextID] = task
	d.nextID++
	return task.ID
}

func (d *db) GetTask(id int) (Task, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	task, exists := d.tasks[id]
	return task, exists
}

func (d *db) UpdateTask(id int, task Task) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.tasks[id]; !exists {
		return false
	}
	task.ID = id
	d.tasks[id] = task
	return true
}

func (d *db) DeleteTask(id int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.tasks[id]; !exists {
		return false
	}
	delete(d.tasks, id)
	return true
}
