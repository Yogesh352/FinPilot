package repository

import (
	"database/sql"
	"stock-api/internal/util"
	"log"
	"fmt"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
}


func (r *UserRepository) RegisterUser(username, password string) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
	`
	hashed, err := util.HashPassword(password)
	if err != nil{
		return err
	}

	_, dbErr := r.db.Exec(query,
		username, hashed,
	)
	if dbErr != nil {
		log.Printf("Failed to insert user: %s", username)
		return err
	}
	

	return nil
}

func (r *UserRepository) LoginUser(username, password string) (string, int, error) {
	var (
		storedHash string
		userID     int
	)

	query := `SELECT id, password FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&userID, &storedHash)
	if err != nil {
		return "", 0, fmt.Errorf("database error: %w", err)
	}

	if !util.CheckPasswordHash(password, storedHash) {
		return "", 0, fmt.Errorf("invalid password")
	}

	token, err := util.GenerateToken(username, userID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, userID, nil
}