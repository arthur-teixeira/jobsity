package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"jobsity-backend/repository"
	"jobsity-backend/service"
	"jobsity-backend/validators"
	"log"
	"net/http"
)

type AuthController struct {
	hasherService *service.Argon2Hasher
	userRepo      *repository.UserRepository
}

func NewAuthController(db *sql.DB, saltLen, keyLen uint32, userRepo *repository.UserRepository) *AuthController {
	hasher := service.NewHasher(saltLen, keyLen)
	return &AuthController{
		hasher,
		userRepo,
	}
}

func (controller AuthController) CreateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", ", POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := validators.ValidateAuthRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := controller.userRepo.GetUserByEmail(body.Email)
	if err != nil {
		log.Println("Error validating existing user: ", err)
		errorResponse(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	if user != nil {
		errorResponse(w, errors.New("A user with this email already exists"), http.StatusBadRequest)
		return
	}

	hashSalt, err := controller.hasherService.GenerateHash([]byte(body.Password), []byte{})
	if err != nil {
		log.Println("Error generating password hash: ", err)
		errorResponse(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	err = controller.userRepo.CreateUser(body.Email, hashSalt)
	if err != nil {
		log.Println("Error saving user to database: ", err)
		errorResponse(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	token, err := service.CreateJWTToken(body.Email)
	if err != nil {
		errorResponse(w, errors.New("An unexpected error occurred"), http.StatusInternalServerError)
		return
	}

	res := AuthResponse{token}
	response, _ := json.Marshal(res)
	okResponse(w, response)
}

func (controller AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", ", POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := validators.ValidateAuthRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := controller.userRepo.GetUserByEmail(body.Email)
	if err != nil {
		log.Println("Error getting user data from database: ", err)
		errorResponse(w, errors.New("An unexpected error occurred"), http.StatusInternalServerError)
		return
	}

	if user == nil {
		errorResponse(w, errors.New("Invalid email or password"), http.StatusUnauthorized)
		return
	}

	isValid := controller.hasherService.Compare(user.Password, user.Salt, []byte(body.Password))
	if !isValid {
		errorResponse(w, errors.New("Invalid email or password"), http.StatusUnauthorized)
		return
	}

	token, err := service.CreateJWTToken(user.Email)
	if err != nil {
		errorResponse(w, errors.New("An unexpected error occurred"), http.StatusInternalServerError)
		return
	}

	res := AuthResponse{token}
	response, _ := json.Marshal(res)
	okResponse(w, response)
}
