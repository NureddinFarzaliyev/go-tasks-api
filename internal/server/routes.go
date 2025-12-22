package server

import (
	"github.com/NureddinFarzaliyev/go-tasks-api/internal/task"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		tasksRoutes(r)
	})
}

func tasksRoutes(r chi.Router) {
	memoryRepo := task.NewInMemoryTaskRepo()
	h := task.NewTaskHandler(memoryRepo)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Edit)
		r.Delete("/{id}", h.Delete)
	})
}
