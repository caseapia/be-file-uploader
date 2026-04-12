package user

import "be-file-uploader/internal/repository/mysql"

type Service struct {
	repo *mysql.Repository
}

func NewService(db *mysql.Repository) *Service {
	return &Service{repo: db}
}
