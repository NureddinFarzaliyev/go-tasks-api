package server

import (
	"database/sql"

	"github.com/NureddinFarzaliyev/go-tasks-api/internal/task"
	"github.com/go-chi/chi/v5"
)

func Routes(r chi.Router, db *sql.DB) {
	r.Route("/v1", func(r chi.Router) {
		tasksMemoryRepo := task.NewInMemoryTaskRepo()
		tasksRoutes(r, tasksMemoryRepo)
	})

	r.Route("/v2", func(r chi.Router) {
		tasksSQLiteRepo := task.NewSQLiteTaskRepo(db)
		tasksRoutes(r, tasksSQLiteRepo)
	})
}

func tasksRoutes(r chi.Router, repo task.TaskRepository) {
	h := task.NewTaskHandler(repo)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Edit)
		r.Delete("/{id}", h.Delete)
	})
}
