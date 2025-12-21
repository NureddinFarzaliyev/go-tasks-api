package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NureddinFarzaliyev/go-tasks-api/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type UpdateTaskRequest struct {
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type TaskHandler struct {
	tasks  []Task
	nextID int
	mu     sync.Mutex
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	httpx.JSON(w, http.StatusOK, httpx.Envelope{"data": h.tasks})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		httpx.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(t.Description) == "" {
		httpx.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	h.nextID++
	t.ID = h.nextID
	t.Completed = false
	t.CreatedAt = time.Now()
	t.CompletedAt = nil

	h.mu.Lock()
	defer h.mu.Unlock()

	h.tasks = append(h.tasks, t)
	httpx.JSON(w, http.StatusCreated, httpx.Envelope{"data": t})
}

func (h *TaskHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)

	var body UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		httpx.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var t *Task

	h.mu.Lock()
	defer h.mu.Unlock()

	for idx, val := range h.tasks {
		if val.ID == idInt {
			t = &h.tasks[idx]
		}
	}

	if t == nil {
		httpx.Error(w, "Task not found", http.StatusNotFound)
		return
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

	httpx.JSON(w, http.StatusOK, httpx.Envelope{"data": t})
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(strId)

	if err != nil {
		httpx.Error(w, "ID is not correct", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var idxToDelete int = -1

	for idx, key := range h.tasks {
		if key.ID == intId {
			idxToDelete = idx
		}
	}

	if idxToDelete == -1 {
		httpx.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	h.tasks = append(h.tasks[:idxToDelete], h.tasks[idxToDelete+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
