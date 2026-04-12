package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) CreateImage(ctx context.Context, tx bun.IDB, image *models.Image) error {
	_, err := tx.NewInsert().
		Model(image).
		Exec(ctx)
	return err
}

func (r *Repository) SearchImageByID(ctx context.Context, id int) (*models.Image, error) {
	image := new(models.Image)
	err := r.DB.NewSelect().
		Model(image).
		Relation("Uploader").
		Where("i.id = ?", id).
		Limit(1).
		Scan(ctx)
	return image, err
}

func (r *Repository) SearchImagesByUserID(ctx context.Context, userID int) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
		Where("i.uploaded_by = ?", userID).
		OrderExpr("i.id DESC").
		Scan(ctx)

	return images, err
}

func (r *Repository) SearchAllImages(ctx context.Context) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.NewSelect().
		Model(&images).
		Relation("Uploader").
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
