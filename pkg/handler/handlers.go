package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
	r.ParseForm()
	c, err := r.Cookie("session")
	if err == nil {
		log.Println(c.Value, "dfdfs")
		fmt.Fprintln(w, "hello")
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
	cookie := &http.Cookie{
		Name:    "session",
		Value:   sessionId,
		Expires: time.Now().AddDate(0, 1, 0),
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if service.RegisterUser(h.DB, login, password) {
		log.Println("create new user", login)
		http.Redirect(w, r, "/main", http.StatusFound)
	} else {
		fmt.Fprintln(w, "Такой логин уже есть")
	}
}
