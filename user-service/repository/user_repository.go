package repository

import (
	"database/sql"
	"errors"
	"github.com/Fagan04/Penguin-Chat-App/user-service/models"
	_ "github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	DB *sql.DB
}

func (repo *UserRepository) CreateUser(user models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)", user.ID, user.Username, user.Email, user.Password)
	return err
}

func (repo *UserRepository) GetUserBYID(id int) (models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user for this email address not found")
	}
	return user, nil
}