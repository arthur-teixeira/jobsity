package controller

import (
	"database/sql"
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

func NewAuthController(db *sql.DB, saltLen, keyLen uint32) *AuthController {
	hasher := service.NewHasher(saltLen, keyLen)
	userRepo := repository.NewUserRepository(db)

	return &AuthController{
		hasher,
		userRepo,
	}
}

func (controller AuthController) CreateUser(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := validators.ValidateCreateUserRequest(r)
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

	err = controller.userRepo.CreateUser(body.Name, body.Email, hashSalt)
	if err != nil {
		log.Println("Error saving user to database: ", err)
		errorResponse(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	okCreated(w)
}

func (controller AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
  enableCors(&w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := validators.ValidateSignInRequest(r)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := controller.userRepo.GetUserByEmail(body.Email)
	if err != nil {
		log.Println("Error getting user data from database: ", err)
		errorResponse(w, errors.New("Unexpected error"), http.StatusInternalServerError)
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

	noContent(w)
}
