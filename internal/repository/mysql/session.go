package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) CreateSession(ctx context.Context, tx bun.IDB, session *models.Session) error {
	_, err := tx.NewInsert().
		Model(session).
		Exec(ctx)

	return err
}

func (r *Repository) SearchSessionByRefreshHash(ctx context.Context, hash string) (*models.Session, error) {
	session := new(models.Session)

	err := r.DB.NewSelect().
		Model(session).
		Where("refresh_hash = ?", hash).
		Scan(ctx)

	return session, err
}

func (r *Repository) SearchSessionByID(ctx context.Context, id string) (*models.Session, error) {
	session := new(models.Session)

	err := r.DB.NewSelect().
		Model(session).
		Where("id = ?", id).
		Scan(ctx)

	return session, err
}

func (r *Repository) TerminateSession(ctx context.Context, tx bun.IDB, sessionID string) (bool, error) {
	session := new(models.Session)

	_, err := tx.NewDelete().
		Model(session).
		Where("id = ?", sessionID).
		Exec(ctx)

	return true, err
}
