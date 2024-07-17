package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"testtask/handlers"
	"time"

	"github.com/gorilla/mux"
)

var rMux = mux.NewRouter()

var port = ":8010"

func main() {
	arguments := os.Args
	if len(arguments) >= 2 {
		port = ":" + arguments[1]
	}

	s := http.Server{
		Addr:         port,
		Handler:      rMux,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	rMux.NotFoundHandler = http.HandlerFunc(handlers.DefaultHandler)

	NotAllowed := handlers.NotAllowedHandler{}
	rMux.MethodNotAllowedHandler = NotAllowed

	//GET
	getMux := rMux.Methods(http.MethodGet).Subrouter()
	getMux.HandleFunc("/getall", handlers.GetAllHandler)
	getMux.HandleFunc("/message/{id:[0-9]+}", handlers.GetMessageDataHandler)

	//POST
	postMux := rMux.Methods(http.MethodPost).Subrouter()
	postMux.HandleFunc("/add", handlers.AddHandler)

	//DELETE
	deleteMux := rMux.Methods(http.MethodDelete).Subrouter()
	deleteMux.HandleFunc("/message/{id:[0-9]+}", handlers.DeleteHandler)

	go func() {
		log.Println("Listening to", port)
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	log.Println("Quitting after signal:", sig)
	time.Sleep(5 * time.Second)
	s.Shutdown(nil)
}
