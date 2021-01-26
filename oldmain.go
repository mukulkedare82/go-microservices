package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func mainold() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("Hello World")
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			//rw.Write([]byte("Ooops!"))
			http.Error(rw, "Ooops", http.StatusBadRequest)
			return
		}
		log.Printf("Data:  %s", data)
		log.Println("")

		fmt.Fprintf(rw, "Hello %s!", data)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}
