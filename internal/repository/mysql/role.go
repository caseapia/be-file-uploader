package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) LookupAllRoles(ctx context.Context) (roles []models.Role, err error) {
	roles = make([]models.Role, 0)

	err = r.DB.NewSelect().
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

func (r *Repository) LookupRoleByID(ctx context.Context, id int) (role *models.Role, err error) {
	role = new(models.Role)

	err = r.DB.NewSelect().
		Model(role).
		Where("id = ?", id).
		Scan(ctx)

	return role, err
}

func (r *Repository) CreateRole(ctx context.Context, tx bun.IDB, role models.Role) (createdRole *models.Role, err error) {
	_, err = tx.NewInsert().
		Model(&role).
		Exec(ctx)

	return &role, err
}

func (r *Repository) DeleteRole(ctx context.Context, tx bun.IDB, id int) (err error) {
	role := new(models.Role)

	_, err = tx.NewDelete().
		Model(role).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *Repository) EditRole(ctx context.Context, tx bun.IDB, role models.Role) (err error) {
	_, err = tx.NewUpdate().
		Model(&role).
		Where("id = ?", role.ID).
		Exec(ctx)
	return err
}
