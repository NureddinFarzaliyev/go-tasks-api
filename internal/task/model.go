package task

import "time"

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
