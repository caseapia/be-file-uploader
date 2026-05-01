package roadmap

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/notification"
)

type Service struct {
	repo   *mysql.Repository
	notify *notification.Service
}

func NewService(repo *mysql.Repository, notify *notification.Service) *Service {
	return &Service{repo: repo, notify: notify}
}
