package image

import (
	"be-file-uploader/internal/repository/mysql"
	r2 "be-file-uploader/pkg/storage"
)

type Service struct {
	repo    *mysql.Repository
	storage *r2.Storage
}

func NewService(repo *mysql.Repository, storage *r2.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}
