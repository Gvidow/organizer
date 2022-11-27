package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gvidow/organizer/pkg/service"
)

type task struct {
	SessionId string
	Title     string
	Desk      string
	Year      int
	Month     int
	Day       int
	Hour      int
	Min       int
	Subj      string
	Exam      string
	Mark      int
}
type data struct {
	Data []task
}
type user struct {
	Login    string
	Password string
}

func (h *Handler) apiSignUp(w http.ResponseWriter, r *http.Request) {
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	userReg := &user{}
	err := bodyJSON.Decode(userReg)
	if err != nil {
		fmt.Fprintln(w, `{"status": "error", "info": "body reading error"`)
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

func (h *Handler) apiSignIn(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		log.Println(12)
		fmt.Fprintln(w, `{"status": "ok", "session": "`+cookie.Value+`"}`)
		return
	}
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	userReg := &user{}
	err = bodyJSON.Decode(userReg)
	if err != nil {
		fmt.Fprintln(w, `{"status": "error", "info": "body reading error"`)
		return
	}
	sessionId, err := service.SignIn(h.DB, userReg.Login, userReg.Password)
	if err != nil {
		fmt.Fprintln(w, `{"status": "error", "info": "login failed"`)
		return
	}
	log.Println("handler: signIn: create sessionId", sessionId)
	cookie = &http.Cookie{
		Name:    "session",
		Value:   sessionId,
		Expires: time.Now().AddDate(0, 1, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	log.Println("set cookie")
	fmt.Fprintln(w, `{"status": "ok", "session": "`+sessionId+`"}`)
}

func (h *Handler) apiAddTasks(w http.ResponseWriter, r *http.Request) {
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	tasksAdd := &data{}
	err := bodyJSON.Decode(&tasksAdd)
	if err != nil {
		fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
		return
	}
	log.Println(tasksAdd)
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

func (h *Handler) apiLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	log.Println("handler: apiLogout:", r.Method, err)
	if err == http.ErrNoCookie {
		fmt.Fprintf(w, `{"status": "error", "info": "error cookie"}`)
		return
	}
	if err != nil {
		fmt.Fprintf(w, `{"status": "error", "info": "error cookie"}`)
	}
	cookie.Expires = time.Now().AddDate(0, -1, 0)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}
