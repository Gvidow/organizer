package main

import (
	"log"
	"net/http"

	"github.com/gvidow/organizer/pkg/handler"
)

func main() {
	server := handler.Server{
		Server: &http.Server{Addr: "localhost:8080"},
	}
	err := server.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	server.HangHandlers()
	log.Println("Start")
	log.Fatal(server.Run())
}
