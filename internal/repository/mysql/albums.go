package mysql

import (
	"context"

	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

func (r *Repository) CreateAlbum(ctx context.Context, tx bun.IDB, album models.Album) (*models.Album, error) {
	_, err := tx.NewInsert().
		Model(&album).
		Exec(ctx)
	return &album, err
}

func (r *Repository) LookupAlbumByID(ctx context.Context, id int) (*models.Album, error) {
	album := new(models.Album)

	err := r.DB.NewSelect().
		Model(album).
		Where("al.id = ?", id).
		Relation("CreatedBy").
		Relation("Items").
		Relation("Items.Uploader").
		Limit(1).
		Scan(ctx)

	return album, err
}

func (r *Repository) LookupAllAlbums(ctx context.Context) ([]models.Album, error) {
	albums := make([]models.Album, 0)

	err := r.DB.NewSelect().
		Model(&albums).
		Relation("CreatedBy").
		Relation("Items").
		Relation("Items.Uploader").
		Scan(ctx)
	return albums, err
}

func (r *Repository) DeleteAlbum(ctx context.Context, tx bun.IDB, album *models.Album) error {
	_, err := tx.NewDelete().
		Model(album).
		Where("al.id = ?", album.ID).
		Exec(ctx)

	return err
}
