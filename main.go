package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mukulkedare/go-microservice-tuts/handlers"
)

func main() {

	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create request handlers
	hh := handlers.NewHello(logger)
	gh := handlers.NewGoodBye(logger)

	// create your serveMux instance and register handlers
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	http.ListenAndServe(":9090", sm)
}
