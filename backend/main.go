package main

import (
	"jobsity-backend/controller"
	"jobsity-backend/db"
	"jobsity-backend/middlewares"
	"jobsity-backend/repository"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const (
	SALT_LEN = 1024
	KEY_LEN  = 1024
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

	userRepo := repository.NewUserRepository(db)
	authController := controller.NewAuthController(db, SALT_LEN, KEY_LEN, userRepo)

	http.HandleFunc("/api/tasks", middlewares.AuthMiddleware(userRepo, taskController.GetTasks))
	http.HandleFunc("/api/task", middlewares.AuthMiddleware(userRepo, taskController.HandleTask))
	http.HandleFunc("/api/auth/signup", authController.CreateUser)
	http.HandleFunc("/api/auth/login", authController.SignIn)

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalln("Unable to start server:", err)
	}
}
