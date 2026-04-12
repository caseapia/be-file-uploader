package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) UseInvite(ctx context.Context, tx bun.Tx, invite models.Invite) error {
	_, err := tx.NewUpdate().
		Model(&invite).
		WherePK().
		Set("is_active = ?", false).
		Set("used_by = ?", invite.UsedBy).
		Exec(ctx)
	return err
}

func (r *Repository) SearchAllInvites(ctx context.Context) ([]models.Invite, error) {
	var invites []models.Invite

	err := r.DB.NewSelect().
		Model(&invites).
		Relation("Creator").
		Relation("User").
		Scan(ctx)

	return invites, err
}

func (r *Repository) SearchInviteByCode(ctx context.Context, code string) (*models.Invite, error) {
	invite := new(models.Invite)
	err := r.DB.NewSelect().
		Model(invite).
		Relation("Creator").
		Relation("User").
		Where("code = ?", code).
		Scan(ctx)

	return invite, err
}

func (r *Repository) SearchInviteByID(ctx context.Context, id int) (*models.Invite, error) {
	invite := new(models.Invite)
	err := r.DB.NewSelect().
		Model(invite).
		Relation("Creator").
		Relation("User").
		Where("id = ?", id).
		Scan(ctx)

	return invite, err
}

func (r *Repository) CreateInvite(ctx context.Context, tx bun.IDB, invite *models.Invite) error {
	_, err := tx.NewInsert().
		Model(invite).
		Exec(ctx)
	return err
}

func (r *Repository) RevokeInvite(ctx context.Context, tx bun.Tx, invite *models.Invite) error {
	_, err := tx.NewUpdate().
		Model(invite).
		Set("is_active = ?", false).
		Exec(ctx)

	return err
}
