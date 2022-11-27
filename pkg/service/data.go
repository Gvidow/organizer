package service

import (
	"database/sql"
	"log"

	"github.com/gvidow/organizer/pkg/repository"
)

type Data struct {
	Status   string            `json:"status,omitempty"`
	Data     []repository.Task `json:"data"`
	UserName string            `json:"username,omitempty"`
}

func AddTasks(db *sql.DB, data *Data) (int, error) {
	for i, row := range data.Data {
		userId, _, err := repository.GetUser(db, row.SessionId)
		if err != nil {
			return i, err
		}
		err = repository.AddTask(db, userId, &row)
		log.Println(userId, err)
		if err != nil {
			return i, err
		}
	}
	return 0, nil
}

func GetTasksAll(db *sql.DB, sessionId string) (*Data, error) {
	userId, userLogin, err := repository.GetUser(db, sessionId)
	if err != nil {
		log.Println("service: getTasksAll:", err)
		return nil, err
	}

	dataTask, err := repository.UserTasksAll(db, userId, sessionId)
	if err != nil {
		return nil, err
	}
	return &Data{Status: "ok", Data: dataTask, UserName: userLogin}, nil

}
