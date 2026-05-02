package auth

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/pkg/geo"
)

type Service struct {
	repo *mysql.Repository
	geo  *geo.Service
}

func NewService(db *mysql.Repository, geo *geo.Service) *Service {
	return &Service{repo: db, geo: geo}
}
