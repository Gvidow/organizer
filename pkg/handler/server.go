package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server
	customHandler Handler
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}

func (s *Server) InitDB() error {
	return s.customHandler.Service.ConnectDB()
}

func (s *Server) HangHandlers() {
	//r := http.NewServeMux()
	//r.HandleFunc("/", redirectToMain)
	router := mux.NewRouter()
	router.HandleFunc("/", s.customHandler.redirectToMain)
	router.HandleFunc("/main", s.customHandler.mainPage).Methods("GET")
	router.HandleFunc("/updatecookie", s.customHandler.updateCookie).Methods("POST")
	router.HandleFunc("/user/sign{sign}", s.customHandler.writeSignIn).Methods("GET")
	router.HandleFunc("/user/signin", s.customHandler.signIn).Methods("POST")
	router.HandleFunc("/user/signup", s.customHandler.signUp).Methods("POST")
	//router.HandleFunc("/user/home", s.customHandler.userHome).Methods("GET")
	//router.HandleFunc("/office", s.customHandler.office).Methods("GET")
	s.Handler = router
}
