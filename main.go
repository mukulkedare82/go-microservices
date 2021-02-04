package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/mukulkedare/go-microservice-tuts/handlers"
)

func main() {

	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create request handlers
	ph := handlers.NewProducts(logger)

	// create your serveMux instance and register handlers
	sm := mux.NewRouter() // using gorilla mux router

	// register specific handlers with gorilla mux
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)

	// create http server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// listen is blocking call, so wrap it in a go func/routine
	// this will not block in main and we can handle shutdown next
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// register for signal
	// create a channel to notify
	// signal broadcasts to this channel whenever sigterm or any interrupt occurs
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan // blocks here
	logger.Println("Received terminate, graceful shutdown", sig)

	// gracefull shutdown
	// close clients gracefully, wait for all incomplete requests to finish
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
