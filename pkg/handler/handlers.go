package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gvidow/organizer/pkg/service"
)

type Handler struct {
	service.Service
}

func (h *Handler) redirectToMain(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main", http.StatusFound)
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil {
		log.Println("handler: mainPage: user in ", c.Value, "session")
		WriteFile(w, PathFront+"html/main_user.html")
		return
	}
	if err := WriteFile(w, PathFront+"html/main_page.html"); err != nil {
		log.Println(err)
		WriteError(w)
	}
}

func (h *Handler) writeSignIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch vars["sign"] {
	case "in":
		WriteFile(w, PathFront+"html/signin.html")
	case "up":
		WriteFile(w, PathFront+"html/signup.html")
	default:
		http.Redirect(w, r, "/main", http.StatusFound)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	log.Println(login, password)
	sessionId, err := service.SignIn(h.DB, login, password)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	log.Println("handler: signIn: create sessionId", sessionId)
	cookie := &http.Cookie{
		Name:    "session",
		Value:   sessionId,
		Expires: time.Now().AddDate(0, 1, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	log.Println("set cookie")
	http.Redirect(w, r, "/main", http.StatusFound)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	err := service.RegisterUser(h.DB, login, password, true)
	if err == nil {
		log.Println("create new user", login)
		http.Redirect(w, r, "/main", http.StatusFound)
		return
	}
	switch e := err.(*mysql.MySQLError); e.Number {
	case uint16(1062):
		log.Println("handelr: signUp: login exist", login)
		fmt.Fprintln(w, "handelr: signUp: login exist")
	case uint16(1366):
		log.Println("handelr: signUp: error encoding", login)
		fmt.Fprintln(w, "handelr: signUp: error encoding")
	default:
		log.Println("handelr: signUp: new error", login)
		fmt.Fprintln(w, "handelr: signUp: new error")
	}
}

func (h *Handler) updateCookie(w http.ResponseWriter, r *http.Request) {
	newDate := r.FormValue("date")
	c, _ := r.Cookie("session")
	sessionId := c.Value
	log.Println(newDate, sessionId)
	//service.UpdateDate(h.DB, sessionId, newDate)

	//http.Redirect(w, r, "/main", http.StatusFound)
	fmt.Fprintln(w, "New date", newDate)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	log.Println("handler: apiLogout:", r.Method, err)
	if err != nil {
		fmt.Fprintf(w, `{"status": "error", "info": "error cookie"}`)
	}
	cookie.Expires = time.Now().AddDate(0, -1, 0)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/main", http.StatusFound)
}
