package validators

import (
	"encoding/json"
	"errors"
	"io"
	"jobsity-backend/entitites"
	"log"
	"net/http"
	"strconv"
)

func ValidateCreateTaskRequest(r *http.Request) (*TaskRequest, error) {
  var body TaskRequest

  bodyBytes, err := io.ReadAll(r.Body)
  if err != nil { 
    log.Println("Error reading request body: ", err)
    return nil, errors.New("Invalid request body")
  }

  err = json.Unmarshal(bodyBytes, &body)
  if err != nil {
    log.Println("Error parsing JSON body: ", err)
    return nil, errors.New("Invalid request body")
  }

  if body.Title == "" {
    return nil, errors.New("'Title' cannot be empty")
  }

  return &body, nil
}

func ValidateEditTaskRequest(r *http.Request) (*entitites.Task, error) {
  var body entitites.Task

  bodyBytes, err := io.ReadAll(r.Body)
  if err != nil {
    log.Println("Error reading request body: ", err)
    return nil, errors.New("Invalid request body")
  }

  err = json.Unmarshal(bodyBytes, &body)
  if err != nil {
    log.Println("Error parsing JSON body: ", err)
    return nil, errors.New("Invalid request body")
  }

  if body.Title == "" {
    return nil, errors.New("'title' cannot be empty")
  }

  if body.Id == 0 {
    return nil, errors.New("'id' cannot be empty")
  }

  return &body, nil
}

func ValidateDeleteTaskRequest(r *http.Request) (int, error) {
  taskId := r.URL.Query().Get("id") 
  if taskId == "" {
    return -1, errors.New("'id' cannot be empty")
  }

  taskIdInt, err := strconv.Atoi(taskId)
  if err != nil {
    return -1, errors.New("'id' must be a valid integer")
  }

  if taskIdInt == 0 {
    return -1, errors.New("'id' must be greater than zero")
  }

  return taskIdInt, nil
}
