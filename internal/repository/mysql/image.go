package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateImage(ctx context.Context, tx bun.IDB, image *models.Image) error {
	_, err := tx.NewInsert().
		Model(image).
		Exec(ctx)
	return err
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

func (r *Repository) SearchImageByID(ctx context.Context, id int) (*models.Image, error) {
	image := new(models.Image)

	err := r.DB.NewSelect().
		Model(image).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Where("i.id = ?", id).
		Limit(1).
		Scan(ctx)

	return image, err
}

func (r *Repository) SearchOwnImages(ctx context.Context, user *models.User) ([]models.Image, error) {
	var images []models.Image

	err := r.DB.NewSelect().
		Model(&images).
		Where("uploaded_by = ?", user.ID).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		OrderExpr("i.id DESC").
		Scan(ctx)
	return images, err
}

func (r *Repository) SearchImagesByUserID(ctx context.Context, userID int) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		Where("i.uploaded_by = ?", userID).
		Where("i.is_private = ?", false).
		OrderExpr("i.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) SearchAllImages(ctx context.Context) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
		Relation("Album").
		Relation("Album.CreatedBy").
		Relation("Likes").
		Relation("Likes.Author").
		OrderExpr("i.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) DeleteImage(ctx context.Context, tx bun.IDB, image *models.Image) error {
	_, err := tx.NewDelete().
		Model(image).
		WherePK().
		Exec(ctx)

	return err
}

func (r *Repository) UpdateImage(ctx context.Context, tx bun.IDB, image *models.Image) (*models.Image, error) {
	_, err := tx.NewUpdate().
		Model(image).
		WherePK().
		Exec(ctx)

	return image, err
}

// func (r *Repository) AddView(ctx context.Context, tx bun.IDB, views models.ImageViews) (inserted bool, err error) {
// 	res, err := tx.NewInsert().
// 		Model(&views).
// 		On("CONFLICT (image_id, author_id) DO NOTHING").
// 		Exec(ctx)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	rows, _ := res.RowsAffected()
// 	return rows > 0, nil
// }

func (r *Repository) AddComment(ctx context.Context, tx bun.IDB, comment models.ImageComments) (inserted bool, err error) {
	res, err := tx.NewInsert().
		Model(&comment).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *Repository) AddLike(ctx context.Context, tx bun.IDB, like models.ImageLikes) (inserted bool, err error) {
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

func (r *Repository) RemoveLike(ctx context.Context, tx bun.IDB, like models.ImageLikes) (delete bool, err error) {
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
