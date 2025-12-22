package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/NureddinFarzaliyev/go-tasks-api/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	repo TaskRepository
}

func NewTaskHandler(repo TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	tasks := h.repo.Get()
	httpx.JSON(w, http.StatusOK, httpx.Envelope{"data": tasks})
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

	data, err := h.repo.Create(t)

	if err != nil {
		httpx.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.JSON(w, http.StatusCreated, httpx.Envelope{"data": data})
}

func (h *TaskHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)

	var body UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		switch err {
		case ErrTaskNotFound:
			httpx.Error(w, err.Error(), http.StatusNotFound)
		default:
			httpx.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	t, err := h.repo.Edit(idInt, body)

	if err != nil {
		httpx.Error(w, err.Error(), http.StatusNotFound)
		return
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

	err = h.repo.Delete(intId)

	if err != nil {
		switch err {
		case ErrTaskNotFound:
			httpx.Error(w, err.Error(), http.StatusNotFound)
		default:
			httpx.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
