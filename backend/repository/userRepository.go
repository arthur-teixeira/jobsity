package repository

import (
	"database/sql"
	"jobsity-backend/entitites"
	"jobsity-backend/service"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (repo UserRepository) CreateUser(name, email string, hashSalt *service.HashSalt) error {
	query := "INSERT INTO users (name, email, password, salt) VALUES ($1, $2, $3, $4)"

	_, err := repo.db.Exec(query, name, email, hashSalt.Hash, hashSalt.Salt)
	if err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) GetUserByEmail(email string) (*entitites.User, error) {
	var user entitites.User

	query := "SELECT id, name, email, password, salt FROM users WHERE email = $1"

	row := repo.db.QueryRow(query, email)

  err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Salt)
	if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }

		return nil, err
	}

	return &user, nil
}
