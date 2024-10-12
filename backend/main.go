package main

import (
	"jobsity-backend/controller"
	"jobsity-backend/db"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)

const (
  SALT_LEN = 1024
  KEY_LEN = 1024
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Println(err)
  }

	db, err := db.GetConn()
	if err != nil {
		log.Fatalln(err)
	}

	taskController := controller.NewTaskController(db)
  authController := controller.NewAuthController(db, SALT_LEN, KEY_LEN)

  http.HandleFunc("/api/tasks", taskController.GetTasks)
  http.HandleFunc("/api/task", taskController.HandleTask)
  http.HandleFunc("/api/auth/signup", authController.CreateUser)
  http.HandleFunc("/api/auth/login", authController.SignIn)

  err = http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatalln("Unable to start server:", err)
  }
}
