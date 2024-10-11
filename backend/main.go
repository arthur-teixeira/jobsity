package main

import (
	"jobsity-backend/controller"
	"log"
	"net/http"
)

func main() {
	controller, err := controller.NewController()
	if err != nil {
		log.Fatalln(err)
	}

  http.HandleFunc("/api/tasks", controller.GetTasks)
  http.HandleFunc("/api/task", controller.HandleTask)

  err = http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatalln("Unable to start server:", err)
  }
}
