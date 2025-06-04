package postgres

import (
	"context"
	"database/sql"

	"github.com/VaneZ444/auth-service/internal/app"
	"github.com/VaneZ444/auth-service/internal/repository"
)

type AppRepository struct {
	db *sql.DB
}

func NewAppRepository(db *sql.DB) repository.AppRepository {
	return &AppRepository{db: db}
}

func (r *AppRepository) GetAppByID(ctx context.Context, appID int32) (name string, secret string, err error) {
	const query = `SELECT name, secret FROM apps WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, appID).Scan(&name, &secret)
	if err == sql.ErrNoRows {
		return "", "", app.ErrAppNotFound
	}
	return
}
