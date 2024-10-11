package repository

import (
	"database/sql"
	"jobsity-backend/entitites"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return TaskRepository{
		db,
	}
}

func (r *TaskRepository) GetTasks(searchTerm string) (*[]entitites.Task, error) {
	query := "SELECT id, title, is_completed FROM tasks"
	if searchTerm != "" {
		query += " WHERE ts @@ phraseto_tsquery('english', $1)"
	}
	query += " ORDER BY id ASC"

	var err error
	var rows *sql.Rows
	if searchTerm != "" {
		rows, err = r.db.Query(query, searchTerm)
	} else {
		rows, err = r.db.Query(query)
  }

	if err != nil {
		return nil, err
	}

	tasks := make([]entitites.Task, 0)

	for rows.Next() {
		var current_task entitites.Task

		err = rows.Scan(
			&current_task.Id,
			&current_task.Title,
			&current_task.IsCompleted,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, current_task)
	}

	return &tasks, nil
}

func (r *TaskRepository) CreateTask(title string) (*entitites.Task, error) {
	var task entitites.Task
	task.IsCompleted = false
	task.Title = title

	query := "INSERT INTO tasks (title) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, title).Scan(&task.Id)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateTask(id int, title string, isCompleted bool) (int, error) {
	isCompletedBit := 0
	if isCompleted {
		isCompletedBit = 1
	}

	query := "UPDATE tasks SET title = $1, is_completed = $2 WHERE id = $3"

	result, err := r.db.Exec(query, title, isCompletedBit, id)
	if err != nil {
		return -1, err
	}

	// err only happens here if the db driver does not support the operations,
	// but Postgres does, so we are ignoring it
	affected, _ := result.RowsAffected()
	return int(affected), nil
}

func (r *TaskRepository) DeleteTask(id int) (int, error) {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return -1, err
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}
