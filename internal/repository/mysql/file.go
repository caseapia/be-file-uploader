package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateFile(ctx context.Context, tx bun.IDB, file *models.File) (models.File, error) {
	_, err := tx.NewInsert().
		Model(file).
		Exec(ctx)
	return *file, err
}

func (r *Repository) ReserveDiskSpace(ctx context.Context, tx bun.Tx, user *models.User, size int64) error {
	res, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("used_storage = used_storage + ?", size).
		Where("id = ?", user.ID).
		Where("used_storage + ? <= upload_limit", size).
		Exec(ctx)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "ERR_QUOTA_EXCEEDED")
	}

	return nil
}

func (r *Repository) SearchFileByID(ctx context.Context, id int) (*models.File, error) {
	image := new(models.File)

	err := r.DB.NewSelect().
		Model(image).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Relation("Comments").
		Relation("Comments.Author").
		Relation("Grants").
		Relation("Grants.User").
		Relation("Grants.GrantedBy").
		Where("f.id = ?", id).
		OrderExpr("f.id DESC").
		Limit(1).
		Scan(ctx)

	return image, err
}

func (r *Repository) SearchOwnFiles(ctx context.Context, user *models.User) ([]models.File, error) {
	var images []models.File

	err := r.DB.NewSelect().
		Model(&images).
		ColumnExpr("f.*").
		Where("f.uploaded_by = ?", user.ID).
		WhereOr("EXISTS (SELECT 1 FROM files_grants AS fg WHERE fg.file_id = f.id AND fg.user_id = ?)", user.ID).
		Group("f.id").
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Relation("Comments").
		Relation("Comments.Author").
		Relation("Grants").
		Relation("Grants.User").
		Relation("Grants.GrantedBy").
		OrderExpr("f.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) SearchFilesByUserID(ctx context.Context, userID int) ([]models.File, error) {
	var images []models.File
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Relation("Comments").
		Relation("Comments.Author").
		Relation("Grants").
		Relation("Grants.User").
		Relation("Grants.GrantedBy").
		Where("f.uploaded_by = ?", userID).
		Where("f.is_private = ?", false).
		OrderExpr("f.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) SearchAllFiles(ctx context.Context) ([]models.File, error) {
	var images []models.File
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Relation("Grants").
		Relation("Grants.User").
		Relation("Grants.GrantedBy").
		OrderExpr("f.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) DeleteFile(ctx context.Context, tx bun.IDB, image *models.File) error {
	_, err := tx.NewDelete().
		Model(image).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) UpdateFile(ctx context.Context, tx bun.IDB, image *models.File) (*models.File, error) {
	_, err := tx.NewUpdate().
		Model(image).
		WherePK().
		Exec(ctx)

	return image, err
}

func (r *Repository) AddLike(ctx context.Context, tx bun.IDB, like models.FileLike) (inserted bool, err error) {
	res, err := tx.NewInsert().
		Model(&like).
		Ignore().
		Exec(ctx)
	if err != nil {
		return false, err
	}

	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *Repository) RemoveLike(ctx context.Context, tx bun.IDB, like models.FileLike) (delete bool, err error) {
	res, err := tx.NewDelete().
		Model(&like).
		Where("image_id = ? AND author = ?", like.ImageID, like.AuthorID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *Repository) AddComment(ctx context.Context, tx bun.IDB, comment *models.FileComment) (createdComment *models.FileComment, err error) {
	res, err := tx.NewInsert().
		Model(comment).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_COMMENT_NOT_FOUND")
	}

	return comment, nil
}

func (r *Repository) GrantAccess(ctx context.Context, tx bun.IDB, grants models.FileGrants) (err error) {
	_, err = tx.NewInsert().
		Model(&grants).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RemoveAccess(ctx context.Context, tx bun.IDB, user, file int) (err error) {
	grants := new(models.FileGrants)

	_, err = tx.NewDelete().
		Model(grants).
		Where("user_id = ? AND file_id = ?", user, file).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
