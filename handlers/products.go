package handlers

import (
	"log"
	"net/http"

	"github.com/mukulkedare/go-microservice-tuts/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.getProducts(rw, req)
}

func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
	lp := data.GetProducts()

	/*
		// JSON MARSHAL
		data, err := json.Marshal(lp)
		if err != nil {
			http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
		}
		rw.Write(data)
	*/

	// JSON ENCODER
	// about using json "Encoder" over "Marshal"
	// encoders writes directly to the io.Writer handle (i.e ResponseWriter) and avoids any buffercopy
	// benefits: this is optimization over marshal
	// for larger data objects the buffer copy will be substantial overhead in case of marshal
	// Also encoder is faster than marshal, this has great performance impact on
	// multithreaded concurrent applications.
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to encode json", http.StatusInternalServerError)
	}

}
