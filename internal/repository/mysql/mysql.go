package mysql

import (
	"be-file-uploader/internal/models"

	"github.com/uptrace/bun"
)

type Repository struct {
	DB *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	db.RegisterModel((*models.UserRole)(nil))

	return &Repository{
		DB: db,
	}
}
