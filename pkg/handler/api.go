package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gvidow/organizer/pkg/service"
)

type task struct {
	Title string
	Desk  string
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
	Subj  string
	Exam  string
	Mark  string
}
type user struct {
	Login    string
	Password string
}

func (h *Handler) api(w http.ResponseWriter, r *http.Request) {
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	userReg := &user{}
	err := bodyJSON.Decode(userReg)
	if err != nil {
		fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
		return
	}
	err = service.RegisterUser(h.DB, userReg.Login, userReg.Password, true)
	if err == nil {
		log.Println("create new user", userReg.Login)
		fmt.Fprint(w, `{"status": "ok"}`)
		return
	}
	switch e := err.(*mysql.MySQLError); e.Number {
	case uint16(1062):
		log.Println("handelr: signUp: login exist", userReg.Login)
		fmt.Fprint(w, `{"status": "error", "info": "login is exist"}`)
	case uint16(1366):
		log.Println("handelr: signUp: error encoding")
		fmt.Fprint(w, `{"status": "error", "info": "data cannot be entered"}`)
	default:
		log.Println("handelr: signUp: new error", userReg.Login)
		fmt.Fprint(w, `{"status": "error", "info": "new erro"}`)
	}
}

func (h *Handler) addTask(w http.ResponseWriter, r *http.Request) {
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	tasksAdd := make([]task, 0)
	err := bodyJSON.Decode(&tasksAdd)
	if err != nil {
		fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
		return
	}
	log.Println("handelr: addTask: request", err, tasksAdd)
	// log.Println("handler: add task: body parse ", taskAdd)
	// err = service.RegisterUser(h.DB, userReg.Login, userReg.Password, true)
	// if err == nil {
	// 	log.Println("create new user", userReg.Login)
	// 	fmt.Fprint(w, `{"status": "ok"}`)
	// 	return
	// }
	// switch e := err.(*mysql.MySQLError); e.Number {
	// case uint16(1062):
	// 	log.Println("handelr: signUp: login exist", userReg.Login)
	// 	fmt.Fprint(w, `{"status": "error", "info": "login is exist"}`)
	// case uint16(1366):
	// 	log.Println("handelr: signUp: error encoding")
	// 	fmt.Fprint(w, `{"status": "error", "info": "data cannot be entered"}`)
	// default:
	// 	log.Println("handelr: signUp: new error", userReg.Login)
	// 	fmt.Fprint(w, `{"status": "error", "info": "new erro"}`)
	// }
}
