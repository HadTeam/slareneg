package auth

import (
	"fmt"
	"time"
)

type DatabaseUserRepository struct {
	// db *sql.DB
}

func NewDatabaseUserRepository( /* db *sql.DB */ ) *DatabaseUserRepository {
	return &DatabaseUserRepository{
		// db: db,
	}
}

func (r *DatabaseUserRepository) CreateUser(username, email string, passwordHash, salt []byte) (*User, error) {
	// query := "INSERT INTO users (id, username, email, password_hash, salt, created_at, last_login_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	// result, err := r.db.Exec(query, ...)

	user := &User{
		ID:           generateUserID(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    time.Now(),
		LastLoginAt:  time.Now(),
	}

	return user, nil
}

func (r *DatabaseUserRepository) GetUserByUsername(username string) (*User, error) {
	// query := "SELECT id, username, email, password_hash, salt, created_at, last_login_at FROM users WHERE username = ?"
	// row := r.db.QueryRow(query, username)
	// var user User
	// err := row.Scan(&user.ID, &user.Username, ...)

	return nil, fmt.Errorf("user not found (database implementation)")
}

func (r *DatabaseUserRepository) UpdateLastLogin(userID string) error {
	// query := "UPDATE users SET last_login_at = ? WHERE id = ?"
	// _, err := r.db.Exec(query, time.Now(), userID)

	return nil
}

func (r *DatabaseUserRepository) UserExists(username string) bool {
	// query := "SELECT COUNT(*) FROM users WHERE username = ?"
	// var count int
	// err := r.db.QueryRow(query, username).Scan(&count)
	// return count > 0

	return false
}
