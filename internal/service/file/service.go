package file

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
	"be-file-uploader/internal/service/upload"
	r2 "be-file-uploader/pkg/storage"
)

type Service struct {
	repo    *mysql.Repository
	notify  *notification.Service
	storage *r2.Storage
	upload  *upload.Service
}

func NewService(repo *mysql.Repository, notify *notification.Service, storage *r2.Storage) *Service {
	return &Service{repo: repo, notify: notify, storage: storage, upload: upload.NewService(storage)}
}
