package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/user-service/models"
	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func (repo *UserRepository) CreateUser(user models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	return err
}

func (repo *UserRepository) GetUserBYID(id int) (models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT * FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user for this email address not found")
	}
	return user, nil
}

func (repo *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, errors.New("user for this user not found")
	}
	return user, nil
}

func (repo *UserRepository) UserExists(username, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err := repo.DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
