package task

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type SQLiteTaskRepo struct {
	db *sql.DB
}

var _ TaskRepository = &SQLiteTaskRepo{}

func NewSQLiteTaskRepo(db *sql.DB) *SQLiteTaskRepo {
	return &SQLiteTaskRepo{
		db: db,
	}
}

func (r *SQLiteTaskRepo) Get() []Task {
	rows, err := r.db.Query(`SELECT id, description, completed, created_at, completed_at FROM tasks`)
	if err != nil {
		return []Task{}
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Description, &t.Completed, &t.CreatedAt, &t.CompletedAt); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	return tasks
}

func (r *SQLiteTaskRepo) Create(t Task) (*Task, error) {
	query := "INSERT INTO tasks (description, created_at) values (?, ?)"
	stmt, err := r.db.Prepare(query)

	if err != nil {
		return nil, errors.New("Error while creating")
	}
	defer stmt.Close()

	now := time.Now()
	result, err := stmt.Exec(t.Description, now)
	if err != nil {
		return nil, errors.New("Error while creating")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("Error retrieving last insert id")
	}

	createdTask := &Task{
		ID:          int(id),
		Description: t.Description,
		CreatedAt:   now,
		Completed:   false,
		CompletedAt: nil,
	}

	return createdTask, nil
}

func (r *SQLiteTaskRepo) Edit(id int, body UpdateTaskRequest) (*Task, error) {
	query := "UPDATE tasks SET"
	args := []any{}
	setParts := []string{}

	var updatedDescription string
	var updatedCompleted bool
	var updatedCompletedAt *time.Time

	// Track which fields are updated for the return value
	if body.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *body.Description)
		updatedDescription = *body.Description
	}
	if body.Completed != nil {
		setParts = append(setParts, "completed = ?")
		args = append(args, *body.Completed)
		updatedCompleted = *body.Completed
		if *body.Completed {
			setParts = append(setParts, "completed_at = ?")
			now := time.Now()
			args = append(args, now)
			updatedCompletedAt = &now
		} else {
			setParts = append(setParts, "completed_at = NULL")
			updatedCompletedAt = nil
		}
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	query += " " + strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, id)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrTaskNotFound
	}
	updatedTask := &Task{
		ID:          id,
		Description: updatedDescription,
		Completed:   updatedCompleted,
		CompletedAt: updatedCompletedAt,
	}
	return updatedTask, nil
}

func (r *SQLiteTaskRepo) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
