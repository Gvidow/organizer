package service

import (
	"database/sql"
	"time"

	"github.com/gvidow/organizer/pkg/repository"
)

func RegisterUser(db *sql.DB, login, password string) bool {
	return repository.AddUser(db, login, HashPassword(password))
}

func SignIn(db *sql.DB, login, password string) (string, error) {
	userId, err := repository.CheckUser(db, login, HashPassword(password))
	if err != nil {
		return "", err
	}
	sessionId := NewIdSession()
	date := time.Now().Format("2006-01-02")
	b := repository.AddSession(db, sessionId, userId, date)
	for !b {
		sessionId := NewIdSession()
		b = repository.AddSession(db, sessionId, userId, date)
	}
	return sessionId, nil
}
