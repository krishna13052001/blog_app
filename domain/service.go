package domain

import (
	"blog_app/db"
)

type domainService struct {
	DB db.Service
}

var _ Service = (*domainService)(nil)

func NewDomainService(repo db.Service) Service {
	return &domainService{
		DB: repo,
	}
}
