package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) LookupAllRoles(ctx context.Context) ([]models.Role, error) {
	roles := make([]models.Role, 0)

	err := r.DB.NewSelect().
		Model(&roles).
		Scan(ctx)
	if roles == nil {
		roles = make([]models.Role, 0)
	}

	return roles, err
}

func (r *Repository) UpdateRole(ctx context.Context, tx bun.IDB, role models.Role, columns ...string) (updatedRole *models.Role, err error) {
	query := tx.NewUpdate().
		Model(role).
		WherePK()

	if len(columns) != 0 {
		query.Column(columns...)
	}

	_, err = query.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *Repository) LookupRoleByID(ctx context.Context, id int) (*models.Role, error) {
	role := new(models.Role)

	err := r.DB.NewSelect().
		Model(role).
		Where("id = ?", id).
		Scan(ctx)

	return role, err
}
