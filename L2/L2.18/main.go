package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"L2.18/internal/application/calendar"
	"L2.18/internal/handler"
	"L2.18/internal/infrastructure/repo"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	flag.Parse()

	eventRepo := repo.NewEventRepoImpl()
	service := calendar.NewService(eventRepo)
	h := handler.NewEventHandler(service)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler.LogMiddleware(h.Mux),
	}

	log.Printf("starting server on :%d", *port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
