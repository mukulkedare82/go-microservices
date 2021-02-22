package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
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
	getRouter.HandleFunc("/products", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddleWareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProduct)
	postRouter.Use(ph.MiddleWareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)

	// using redoc openapi go middleware Redoc for swagger documentation UI
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	// For redoc UI to swagger.yaml file
	// provide access to swagger.yaml as get request using build-in go fileserver
	// serving it from current working directory (http app directory)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS - CROSS ORIGIN REQUESTS
	// It is the mechanism of securing cross origin requests and data transfers
	// between browser and servers. Browser by default blocks cross origin requests.
	// Ex:
	// 1. user logs in to ebank.com, it writes a cookie to browser
	// 2. user visits some malicious website evil.com
	// 3. evil.com makes call to ebank.com (malicious logic)
	// 4. browser checks cookies for ebank.com and forwards it in the request (by evil.com)
	//    to ebank.com server
	// 5. evil.com could be doing some malicious activity in the request (like initiate a bank transfer)
	//    it could possibly do it as browser is sending all cookies of ebank.com in request
	//    to the server.
	// 6. To avoid this situation CORS mechanism comes into play, before sending cookies to the server
	//    CORS mechanism checks with the server "which origin requests it allows ?"
	// 7. The server could say only "frontend.ebank.com", in the case evil.com request is rejected
	//    as it is not the approved access control origin.
	// 8. We can use * to give access to all origins, can be used for public API
	// https://medium.com/@baphemot/understanding-cors-18ad6b478e2b

	// To add CORS we can use gorilla middleware CORS
	origins := gohandlers.AllowedOrigins([]string{"http://localhost:3000"})
	//origins := gohandlers.AllowedOrigins([]string{"*"})
	ch := gohandlers.CORS(origins)

	// create http server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
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
