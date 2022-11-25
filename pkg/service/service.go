package service

import (
	"database/sql"

	"github.com/gvidow/organizer/pkg/repository"
)

type Service struct {
	DB *sql.DB
}

func (s *Service) ConnectDB() error {
	var err error
	s.DB, err = repository.ConnectDB()
	return err
}
