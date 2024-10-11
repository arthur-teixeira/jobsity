package controller

import (
	"encoding/json"
	"errors"
	"jobsity-backend/db"
	"jobsity-backend/repository"
	"jobsity-backend/validators"
	"net/http"
)

type Controller struct {
	tasksRepository repository.TaskRepository
}

func NewController() (*Controller, error) {
	db, err := db.GetConn()
	if err != nil {
		return nil, err
	}

	tasksRepository := repository.NewTaskRepository(db)

	return &Controller{
		tasksRepository,
	}, nil
}

func (controller Controller) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (controller Controller) response(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (controller Controller) noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (controller Controller) errorResponse(w http.ResponseWriter, error error, status int) {
	errorResponse := ErrorResponse{
		Error: error.Error(),
	}
	res, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func (controller Controller) GetTasks(w http.ResponseWriter, r *http.Request) {
	controller.enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	searchTerm := r.URL.Query().Get("search")
	tasks, err := controller.tasksRepository.GetTasks(searchTerm)
	if err != nil {
		controller.errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(tasks)
	controller.response(w, response)
}

func (controller Controller) editTask(w http.ResponseWriter, r *http.Request) {
	body, err := validators.ValidateEditTaskRequest(r)
	if err != nil {
		controller.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.UpdateTask(body.Id, body.Title, body.IsCompleted)
	if err != nil {
		controller.errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if numAffected == 0 {
		controller.errorResponse(w, errors.New("Task not found"), http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(body)
	controller.response(w, response)
}

func (controller Controller) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := validators.ValidateDeleteTaskRequest(r)
	if err != nil {
		controller.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	numAffected, err := controller.tasksRepository.DeleteTask(taskId)
	if err != nil {
		controller.errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if numAffected == 0 {
		controller.errorResponse(w, errors.New("Task not found"), http.StatusNotFound)
		return
	}

	controller.noContent(w)
}

func (controller Controller) createTask(w http.ResponseWriter, r *http.Request) {
	body, err := validators.ValidateCreateTaskRequest(r)
	if err != nil {
		controller.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	newTask, err := controller.tasksRepository.CreateTask(body.Title)
	if err != nil {
		controller.errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(newTask)
	controller.response(w, response)
}

func (controller Controller) HandleTask(w http.ResponseWriter, r *http.Request) {
	controller.enableCors(&w)
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
