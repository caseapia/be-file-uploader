package user

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
)

type Service struct {
	repo   *mysql.Repository
	notify *notification.Service
}

func NewService(db *mysql.Repository, notify *notification.Service) *Service {
	return &Service{repo: db, notify: notify}
}
