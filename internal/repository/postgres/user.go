package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VaneZ444/golang-forum/auth-service/internal/entity"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) SaveUser(user *entity.User) (int64, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	var userID int64
	err := r.db.QueryRowContext(context.Background(), query, user.Email, hashPassword(user.Password)).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	return userID, nil
}
