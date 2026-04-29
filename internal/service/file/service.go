package image

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
	r2 "be-file-uploader/pkg/storage"
)

type Service struct {
	repo    *mysql.Repository
	notify  *notification.Service
	storage *r2.Storage
}

func NewService(repo *mysql.Repository, notify *notification.Service, storage *r2.Storage) *Service {
	return &Service{repo: repo, notify: notify, storage: storage}
}
