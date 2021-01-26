package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	return &Hello{logger}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.logger.Println("Hello World")
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		//rw.WriteHeader(http.StatusBadRequest)
		//rw.Write([]byte("Ooops!"))
		http.Error(rw, "Ooops", http.StatusBadRequest)
		return
	}
	h.logger.Printf("Data:  %s", data)
	h.logger.Println("")

	fmt.Fprintf(rw, "Hello %s!", data)
}
