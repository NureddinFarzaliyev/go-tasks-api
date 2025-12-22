package task

import (
	"sync"
	"time"
)

type InMemoryTaskRepo struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

var _ TaskRepository = &InMemoryTaskRepo{}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		tasks:  []Task{},
		nextID: 0,
	}
}

func (r *InMemoryTaskRepo) Get() []Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tasks
}

func (r *InMemoryTaskRepo) Create(t Task) (*Task, error) {
	r.nextID++
	t.ID = r.nextID
	t.Completed = false
	t.CreatedAt = time.Now()
	t.CompletedAt = nil

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks = append(r.tasks, t)
	return &t, nil
}

func (r *InMemoryTaskRepo) Edit(id int, body UpdateTaskRequest) (*Task, error) {
	var t *Task

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, val := range r.tasks {
		if val.ID == id {
			t = &r.tasks[idx]
		}
	}

	if t == nil {
		return nil, ErrTaskNotFound
	}

	if body.Description != nil {
		t.Description = *body.Description
	}

	if body.Completed != nil {
		t.Completed = *body.Completed
		if *body.Completed {
			now := time.Now()
			t.CompletedAt = &now
		} else {
			t.CompletedAt = nil
		}
	}

	return t, nil
}

func (r *InMemoryTaskRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var idxToDelete int = -1

	for idx, key := range r.tasks {
		if key.ID == id {
			idxToDelete = idx
		}
	}

	if idxToDelete == -1 {
		return ErrTaskNotFound
	}

	r.tasks = append(r.tasks[:idxToDelete], r.tasks[idxToDelete+1:]...)
	return nil
}
