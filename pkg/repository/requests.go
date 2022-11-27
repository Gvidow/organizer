package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func AddUser(db *sql.DB, login, hashPassword string) error {
	_, err := db.Exec("INSERT INTO user(user_login, password_hash) VALUES (?, ?);", login, hashPassword)
	log.Println("db: add user", err)
	//log.Println("db: add user: ", err, reflect.TypeOf(err))
	//log.Println("add user", err, login, hashPassword, len(hashPassword), len([]byte(hashPassword)))
	return err
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

func AddTask(db *sql.DB, id int, task *Task) error {
	date := fmt.Sprintf("%d-%d-%d %d:%d", task.Year, task.Month, task.Day, task.Hour, task.Min)
	log.Println("db: add task: date", date)
	_, err := db.Exec(`INSERT INTO task(user_id, task_title, task_description, task_date, task_subject, task_exam, task_mark)
	VALUES (?, ?, ?, ?, ?, ?, ?);`, id, task.Title, task.Desc, date, task.Subj, task.Exam, task.Mark)
	log.Println("db: add task:", err)
	log.Println("db: add:db", task)
	return err
}

func GetUser(db *sql.DB, sessionId string) (int, string, error) {
	var userId int
	var userLogin string
	row := db.QueryRow("SELECT _session.user_id, user_login FROM _session INNER JOIN user ON _session.user_id = user.user_id WHERE session_id = ?;", sessionId)
	if err := row.Err(); err != nil {
		log.Println("db: getUser: select error", err)
	}
	err := row.Scan(&userId, &userLogin)
	if err != nil {
		log.Println("db: getUser: Scan select user", err)
		return 0, "", err
	}

	return userId, userLogin, nil
}

func UserTasksAll(db *sql.DB, userId int, sessionId string) ([]Task, error) {
	rows, err := db.Query("SELECT task_title, task_description, task_date, task_subject, task_exam, task_mark FROM task WHERE user_id = ? ORDER BY task_mark DESC;", userId)
	if err != nil {
		log.Println("rep: userTaskAll: select error:", err)
		return nil, err
	}
	res := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		var date string
		rows.Scan(&task.Title, &task.Desc, &date, &task.Subj, &task.Exam, &task.Mark)
		var y, mon, d, h, min int
		fmt.Sscanf(date, "%d-%d-%d %d:%d", &y, &mon, &d, &h, &min)
		task.Year = y
		task.Month = mon
		task.Day = d
		task.Hour = h
		task.Min = min
		task.SessionId = sessionId
		res = append(res, task)
		log.Println("rep: userTaskAll: select data:", task, "data string =", date)
	}
	rows.Close()
	return res, nil

}

func DelTask(db *sql.DB, userId int) {
	db.Exec("DELETE FROM task WHERE user_id = ?;", userId)
}
