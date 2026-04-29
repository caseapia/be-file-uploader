package notification

import (
	"be-file-uploader/internal/repository/mysql"
)

type Service struct {
	repo *mysql.Repository
}

func NewService(repo *mysql.Repository) *Service {
	return &Service{repo: repo}
}
