package task

type TaskRepository interface {
	Get() []Task
	Create(t Task) (*Task, error)
	Edit(id int, t UpdateTaskRequest) (*Task, error)
	Delete(id int) error
}
