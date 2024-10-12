package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (controller TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	searchTerm := r.URL.Query().Get("search")
	tasks, err := controller.tasksRepository.GetTasks(searchTerm)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(tasks)
	okResponse(w, response)
}

func (controller TaskController) editTask(w http.ResponseWriter, r *http.Request) {
	body, err := validators.ValidateEditTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.UpdateTask(body.Id, body.Title, body.IsCompleted)
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

func (controller TaskController) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := validators.ValidateDeleteTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.DeleteTask(taskId)
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

func (controller TaskController) createTask(w http.ResponseWriter, r *http.Request) {
	body, err := validators.ValidateCreateTaskRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	newTask, err := controller.tasksRepository.CreateTask(body.Title)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(newTask)
	okResponse(w, response)
}

func (controller TaskController) HandleTask(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	switch r.Method {
	case http.MethodPost:
		controller.createTask(w, r)
	case http.MethodPut:
		controller.editTask(w, r)
	case http.MethodDelete:
		controller.deleteTask(w, r)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
