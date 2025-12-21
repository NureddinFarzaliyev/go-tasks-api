package server

import (
	"github.com/NureddinFarzaliyev/go-tasks-api/internal/task"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		tasksRoutes(r, &task.TaskHandler{})
	})
}

func tasksRoutes(r chi.Router, h *task.TaskHandler) {
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Edit)
		r.Delete("/{id}", h.Delete)
	})
}
