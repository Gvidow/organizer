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
	router.HandleFunc("/update/date", s.customHandler.updateCookie).Methods("POST")
	router.HandleFunc("/user/sign{sign}", s.customHandler.writeSignIn).Methods("GET")
	router.HandleFunc("/user/signin", s.customHandler.signIn).Methods("POST")
	router.HandleFunc("/user/signup", s.customHandler.signUp).Methods("POST")
	router.HandleFunc("/api/user/signup", s.customHandler.apiSignUp).Methods("POST")
	router.HandleFunc("/api/user/signin", s.customHandler.apiSignIn).Methods("POST")
	router.HandleFunc("/api/user/addtask", s.customHandler.apiAddTasks).Methods("POST")
	router.HandleFunc("/user/logout", s.customHandler.logout).Methods("POST")
	router.HandleFunc("/api/user/logout", s.customHandler.apiLogout).Methods("POST")
	router.HandleFunc("/api/user/gettask", s.customHandler.getTasks).Methods("POST")
	//router.HandleFunc("/api/user/updatetask", s.customHandler.apiAddTasks).Method("POST")
	//router.HandleFunc("/user/home", s.customHandler.userHome).Methods("GET")
	//router.HandleFunc("/office", s.customHandler.office).Methods("GET")
	s.Handler = router
}
