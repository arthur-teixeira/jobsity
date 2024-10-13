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

func (r *TaskRepository) GetTasksByUser(user *entitites.User) (*[]entitites.Task, error) {
	query := "SELECT id, title, is_completed FROM tasks WHERE user_id = $1"
	query += " ORDER BY id ASC"

  rows, err := r.db.Query(query, user.Id)

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

func (r *TaskRepository) CreateTask(title string, user *entitites.User) (*entitites.Task, error) {
	var task entitites.Task
	task.IsCompleted = false
	task.Title = title

	query := "INSERT INTO tasks (title, user_id) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, title, user.Id).Scan(&task.Id)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateTask(id int, title string, isCompleted bool, user *entitites.User) (int, error) {
	isCompletedBit := 0
	if isCompleted {
		isCompletedBit = 1
	}

	query := "UPDATE tasks SET title = $1, is_completed = $2 WHERE id = $3 AND user_id = $4"

	result, err := r.db.Exec(query, title, isCompletedBit, id, user.Id)
	if err != nil {
		return -1, err
	}

	// err only happens here if the db driver does not support the operations,
	// but Postgres does, so we are ignoring it
	affected, _ := result.RowsAffected()
	return int(affected), nil
}

func (r *TaskRepository) DeleteTask(id int, user *entitites.User) (int, error) {
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	result, err := r.db.Exec(query, id, user.Id)
	if err != nil {
		return -1, err
	}

	affected, _ := result.RowsAffected()
	return int(affected), nil
}
