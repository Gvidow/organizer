package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const sep = "========================================"

func (h *Handler) apiSignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("request", r)
	//io.Copy(os.Stdout, r.Body)
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	bodyJSON := json.NewDecoder(r.Body)
	userReg := &user{}
	err := bodyJSON.Decode(userReg)
	if err != nil {
		log.Println("api sign", err)
		fmt.Fprintln(w, `{"status": "error", "info": "body reading error"`)
		return
	}
	log.Println("aaaaaaaaaaaaaaa", err, "login", userReg.Login, userReg.Password)
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
	log.Println("Idy na vhod")
	cok, err := r.Cookie("session")
	log.Println(err, cok)
	log.Println("request", r)
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
	fmt.Print(sep, "\nREQUEST /api/user/addtask\n\n")
	defer fmt.Println("\nRequest end\n", sep)
	log.Println("api: addTask: ", r)
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	fmt.Println(sep)
	//io.Copy(os.Stdout, r.Body)
	b, err := ioutil.ReadAll(r.Body)
	//str := string(b[1 : len(b)-1])
	str := string(b)
	//str = str[1 : len(str)-1]
	log.Println("0 char", str[0], string(str[len(str)-1]))
	log.Println("[]byte", err, str)
	//      n   bodyJSON := json.NewDecoder(r.Body)
	// tasksAdd := &data{}
	// err := bodyJSON.Decode(&tasksAdd)
	// if err != nil {
	// 	fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
	// 	return
	// }
	tasksAdd := &service.Data{}
	//          err := bodyJSON.Decode(tasksAdd)
	err = json.Unmarshal([]byte(str), tasksAdd)
	if err != nil {
		log.Println("decode error", err)
		fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
		return
	}
	log.Println(tasksAdd)
	n, err := service.AddTasks(h.DB, tasksAdd)
	log.Println(n, err)
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

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Print(sep, "\nREQUEST /api/user/gettask\n\n")
	defer fmt.Println("\nRequest end\n", sep)
	log.Println("api: getTask: ", r)
	fmt.Println(sep)
	session := make(map[string]string)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&session)
	if err != nil {
		log.Println("api: getTask: decode", err)
		fmt.Fprintln(w, `{"status": "error", "info": "error decode"}`)
		return
	}
	sessionId, ok := session["session"]
	if !ok {
		log.Println("api: getTask: get session error")
		fmt.Fprintln(w, `{"status": "error", "info": "get session error"}`)
		return
	}
	res, err := service.GetTasksAll(h.DB, sessionId)
	if err != nil {
		log.Println("api: getTask: get task all error:", err)
		fmt.Fprintln(w, `{"status": "error", "info": "no rows in result set"}`)
		return
	}
	log.Println("api: getTask: get task ok:", res)
	fmt.Println(sep)
	resByte, err := json.Marshal(res)
	if err != nil {
		log.Println("api: getTask: get task marshal:", err)
		fmt.Fprintln(w, `{"status": "error", "info": "ferror data marshal"}`)
		return
	}
	fmt.Fprintln(w, string(resByte))
}

func (h *Handler) apiUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	defer fmt.Print("\n\n\n")
	fmt.Print(sep, "\nREQUEST /api/user/updatetask\n\n")
	defer fmt.Println("\nRequest end\n", sep)
	log.Println("api: updateTask: ", r)
	t := r.Header["Content-Type"]
	log.Println("api", t, r.Method)
	fmt.Println(sep)
	//io.Copy(os.Stdout, r.Body)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api: update: read body error")
	}
	//str := string(b[1 : len(b)-1])
	str := string(b)
	//str = str[1 : len(str)-1]
	// log.Println("0 char", str[0], string(str[len(str)-1]))
	// log.Println("[]byte", err, str)
	//      n   bodyJSON := json.NewDecoder(r.Body)
	// tasksAdd := &data{}
	// err := bodyJSON.Decode(&tasksAdd)
	// if err != nil {
	// 	fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
	// 	return
	// }
	tasksAdd := &service.Data{}
	//          err := bodyJSON.Decode(tasksAdd)
	err = json.Unmarshal([]byte(str), tasksAdd)
	if err != nil {
		log.Println("decode error", err)
		fmt.Fprintln(w, `{"status": "ok", "info": "body reading ok"`)
		return
	}
	if err != nil {
		log.Println("")
	}
	log.Println(tasksAdd)
	userId, sessionId := service.GetID(h.DB, tasksAdd)
	service.UpdateTasks(h.DB, userId, tasksAdd)
	//log.Println(n, err)
	log.Println("handelr: addTask: request", err, tasksAdd)

	////
	res, err := service.GetTasksAll(h.DB, sessionId)
	if err != nil {
		log.Println("api: getTask: get task all error:", err)
		fmt.Fprintln(w, `{"status": "error", "info": "no rows in result set"}`)
		return
	}
	log.Println("api: getTask: get task ok:", res)
	fmt.Println(sep)
	resByte, err := json.Marshal(res)
	if err != nil {
		log.Println("api: getTask: get task marshal:", err)
		fmt.Fprintln(w, `{"status": "error", "info": "ferror data marshal"}`)
		return
	}
	fmt.Fprintln(w, string(resByte))
}
