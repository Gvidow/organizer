package repository

import (
	"database/sql"
	"errors"
	"log"
)

func AddUser(db *sql.DB, login, hashPassword string) bool {
	_, err := db.Exec("INSERT INTO user(user_login, password_hash) VALUES (?, ?);", login, hashPassword)
	log.Println("add user", err, login, hashPassword, len(hashPassword), len([]byte(hashPassword)))
	return err == nil
}

func CheckUser(db *sql.DB, login, hashPassword string) (int, error) {
	row := db.QueryRow("SELECT user_id, user_login, password_hash FROM user WHERE user_login = ?;", login)
	var userId int
	var userLogin, passwordHash string
	err := row.Scan(&userId, &userLogin, &passwordHash)
	log.Println(err)
	if err != nil {
		return 0, err
	}
	if passwordHash != hashPassword {
		return 0, errors.New("неверный логин")
	}
	return userId, nil
}

// func IsNotUseSession(db *sql.DB, sessionId string) bool {
// 	row := db.QueryRow("SELECT session_id, COUNT(*) AS count FROM _session GROUP BY session_id WHERE session_id = ?;", sessionId)
// 	var count int
// 	err := row.Scan(&sessionId, &count)
// 	if err != nil {
// 		return true
// 	}
// 	return count <= 0
// }

func AddSession(db *sql.DB, sessionId string, userId int, date string) bool {
	_, err := db.Exec("INSERT INTO _session(session_id, user_id, use_date) VALUES (?, ?, ?);", sessionId, userId, date)
	log.Println(err, len(sessionId))
	return err == nil
}
