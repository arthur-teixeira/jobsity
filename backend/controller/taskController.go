package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"jobsity-backend/entitites"
	"jobsity-backend/repository"
	"jobsity-backend/validators"
	"net/http"
)

type TaskController struct {
	tasksRepository repository.TaskRepository
}

func NewTaskController(db *sql.DB) *TaskController {
	tasksRepository := repository.NewTaskRepository(db)
	return &TaskController{
		tasksRepository,
	}
}

func (controller TaskController) GetTasks(w http.ResponseWriter, r *http.Request, user *entitites.User) {
  enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := controller.tasksRepository.GetTasksByUser(user)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(tasks)
	okResponse(w, response)
}

func (controller TaskController) editTask(w http.ResponseWriter, r *http.Request, user *entitites.User) {
	body, err := validators.ValidateEditTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.UpdateTask(body.Id, body.Title, body.IsCompleted, user)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if numAffected == 0 {
		errorResponse(w, errors.New("Task not found"), http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(body)
	okResponse(w, response)
}

func (controller TaskController) deleteTask(w http.ResponseWriter, r *http.Request, user *entitites.User) {
	taskId, err := validators.ValidateDeleteTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.DeleteTask(taskId, user)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if numAffected == 0 {
		errorResponse(w, errors.New("Task not found"), http.StatusNotFound)
		return
	}

	noContent(w)
}

func (controller TaskController) createTask(w http.ResponseWriter, r *http.Request, user *entitites.User) {
	body, err := validators.ValidateCreateTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	newTask, err := controller.tasksRepository.CreateTask(body.Title, user)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(newTask)
	okResponse(w, response)
}

func (controller TaskController) HandleTask(w http.ResponseWriter, r *http.Request, user *entitites.User) {
  enableCors(&w)
	switch r.Method {
	case http.MethodPost:
		controller.createTask(w, r, user)
	case http.MethodPut:
		controller.editTask(w, r, user)
	case http.MethodDelete:
		controller.deleteTask(w, r, user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
