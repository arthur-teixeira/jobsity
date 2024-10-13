package middlewares

import (
	"jobsity-backend/entitites"
	"jobsity-backend/repository"
	"jobsity-backend/service"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)
type SecuredHttpHandler func(w http.ResponseWriter, r *http.Request, user *entitites.User)

func AuthMiddleware(userRepo *repository.UserRepository, next SecuredHttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusOK)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := service.ValidateToken(tokenString)
		if token == nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := userRepo.GetUserByEmail(email)
		if err != nil || user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r, user)
	}
}
