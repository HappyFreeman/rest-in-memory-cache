package repo

import (
	"github.com/pkg/errors"
	"sync"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

// SQL-запрос на вставку задачи

type repository struct {
	tasks         map[int]Task
	autoincrement int
	mu            sync.RWMutex
}

// Repository - интерфейс с методом создания задачи
type Repository interface {
	CreateTask(task Task) (int, error)
	GetTask(id int) (Task, error)
	GetTasks() ([]Task, error)
	DeleteTask(id int) error
	UpdateTask(id int, task Task) error
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewRepository() (Repository, error) {

	rep := &repository{
		tasks:         make(map[int]Task),
		autoincrement: 0,
	}

	return rep, nil
}

func (r *repository) CreateTask(task Task) (int, error) {
	var id int

	r.mu.Lock()
	defer r.mu.Unlock()

	r.autoincrement++
	r.tasks[r.autoincrement] = task
	id = r.autoincrement
	return id, nil
}

func (r *repository) GetTask(id int) (Task, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.tasks[id]; !ok {
		return Task{}, errors.New("task not found")
	}

	return r.tasks[id], nil
}

func (r *repository) GetTasks() ([]Task, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *repository) DeleteTask(id int) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)

	return nil
}

func (r *repository) UpdateTask(id int, task Task) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return errors.New("task not found")
	}

	r.tasks[id] = task

	return nil
}
