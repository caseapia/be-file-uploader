package user

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
	"be-file-uploader/internal/service/upload"
	r2 "be-file-uploader/pkg/storage"
)

type Service struct {
	repo   *mysql.Repository
	notify *notification.Service
	upload *upload.Service
}

func NewService(db *mysql.Repository, notify *notification.Service, storage *r2.Storage) *Service {
	return &Service{repo: db, notify: notify, upload: upload.NewService(storage)}
}
