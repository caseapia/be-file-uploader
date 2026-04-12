package mysql

import (
	"context"
	"time"

	"be-file-uploader/internal/models"

	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateSession(ctx context.Context, tx bun.IDB, session *models.Session) error {
	res, err := tx.NewInsert().
		Model(session).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	slog.WithData(slog.M{
		"rows_affected": rows,
		"session_id":    session.ID,
	}).Info("CreateSession rows affected")

	return nil
}

func (r *Repository) SearchSessionByRefreshHash(ctx context.Context, hash string) (*models.Session, error) {
	session := new(models.Session)

	err := r.DB.NewSelect().
		Model(session).
		Where("refresh_hash = ?", hash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *Repository) SearchSessionByID(ctx context.Context, id string) (*models.Session, error) {
	session := new(models.Session)

	err := r.DB.NewSelect().
		Model(session).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)

	return session, err
}

func (r *Repository) CleanupExpiredSessions(ctx context.Context, tx bun.IDB, user *models.User) error {
	_, err := tx.NewDelete().
		Model((*models.Session)(nil)).
		Where("user_id = ?", user.ID).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}

func (r *Repository) TerminateSession(ctx context.Context, tx bun.IDB, sessionID string) (bool, error) {
	_, err := tx.NewDelete().
		Model((*models.Session)(nil)).
		Where("id = ?", sessionID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, err
}

func (r *Repository) UpdateSession(ctx context.Context, tx bun.IDB, s *models.Session) error {
	_, err := tx.NewUpdate().
		Model(s).
		WherePK().
		Exec(ctx)
	return err
}
