package handlers

import (
	"log"
	"net/http"
)

type GoodBye struct {
	logger *log.Logger
}

func NewGoodBye(logger *log.Logger) *GoodBye {
	return &GoodBye{logger}
}

func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Byee"))
}
