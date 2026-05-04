package mysql

import (
	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

type Repository struct {
	DB *bun.DB
	// redis *redis.Client
}

func NewRepository(db *bun.DB) *Repository {
	db.RegisterModel((*models.UserRole)(nil))
	db.RegisterModel((*models.FileGrants)(nil))

	return &Repository{
		DB: db,
		// redis: rdb,
	}
}
