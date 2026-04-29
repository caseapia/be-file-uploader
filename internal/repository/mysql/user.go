package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) LookupUserByName(ctx context.Context, name string) (*models.User, error) {
	user := new(models.User)

	err := r.DB.NewSelect().
		Model(user).
		Where("u.username = ?", name).
		Relation("Roles").
		Relation("Storage").
		Relation("Storage.Uploader").
		Relation("Storage.Likes").
		Relation("Albums").
		Relation("Albums.CreatedBy").
		Relation("Albums.Items").
		Limit(1).
		Scan(ctx)
	return user, err
}

func (r *Repository) LookupUserByID(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)

	err := r.DB.NewSelect().
		Model(user).
		Where("u.id = ?", id).
		Relation("Roles").
		Relation("Storage").
		Relation("Storage.Uploader").
		Relation("Storage.Likes").
		Relation("Albums").
		Relation("Albums.CreatedBy").
		Relation("Albums.Items").
		Limit(1).
		Scan(ctx)

	return user, err
}

func (r *Repository) LookupUsers(ctx context.Context, limit int) ([]models.User, error) {
	users := make([]models.User, 0)

	err := r.DB.NewSelect().
		Model(&users).
		Relation("Roles").
		Limit(limit).
		Scan(ctx)
	return users, err
}

func (r *Repository) UpdateUser(ctx context.Context, tx bun.IDB, user *models.User, columns ...string) (updatedUser *models.User, err error) {
	query := tx.NewUpdate().
		Model(user).
		WherePK()

	if len(columns) > 0 {
		query.Column(columns...)
	} else {
		query.ExcludeColumn("created_at")
	}

	_, err = query.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, tx bun.IDB, user models.User) (*models.User, error) {
	_, err := tx.NewInsert().
		Model(&user).
		Exec(ctx)

	return &user, err
}

func (r *Repository) AddUserInRole(ctx context.Context, tx bun.IDB, userID, roleID int) error {
	_, err := tx.NewInsert().
		Model(&models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}).
		Exec(ctx)

	return err
}

func (r *Repository) RemoveUserFromRole(ctx context.Context, tx bun.IDB, userID, roleID int) error {
	_, err := tx.NewDelete().
		Table("user_roles").
		Where("user_id = ?", userID).
		Where("role_id = ?", roleID).
		Exec(ctx)

	return err
}

func (r *Repository) LookupUserByToken(ctx context.Context, token string) (*models.User, error) {
	user := new(models.User)

	err := r.DB.NewSelect().
		Model(user).
		Where("sharex_token = ?", token).
		Limit(1).
		Scan(ctx)

	return user, err
}
