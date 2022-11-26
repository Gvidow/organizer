package service

import (
	"database/sql"
	"log"

	"github.com/gvidow/organizer/pkg/repository"
)

func AddTasks(db *sql.DB, rows []map[string]string) (int, error) {
	for _, row := range rows {
		k := repository.AddTask(db, row["title"], row["disc"], row["subj"], row["exam"])
		log.Println(k)
	}
	return 0, nil
}
