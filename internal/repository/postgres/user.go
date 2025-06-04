package postgres

import (
	"context"
	"database/sql"

	"github.com/VaneZ444/auth-service/internal/entity"
	"github.com/VaneZ444/auth-service/internal/repository"
	"github.com/VaneZ444/auth-service/internal/usecase"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *entity.User) (int64, error) {
	const query = `INSERT INTO users (email, password, role, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Hash, user.Role, user.Status).Scan(&user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `SELECT id, email, password, role, status FROM users WHERE email = $1`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Hash, &user.Role, &user.Status)
	if err == sql.ErrNoRows {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const query = `SELECT role FROM users WHERE id = $1`
	var role entity.Role
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return false, usecase.ErrUserNotFound
	}
	if err != nil {
		return false, err
	}
	return role == entity.RoleAdmin, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	const query = `SELECT id, email, password, role, status FROM users WHERE id = $1`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Hash, &user.Role, &user.Status)
	if err == sql.ErrNoRows {
		return nil, usecase.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
