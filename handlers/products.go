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
	if req.Method == http.MethodGet {
		p.getProducts(rw, req)
		return
	}

	if req.Method == http.MethodPost {
		p.addProduct(rw, req)
		return
	}

	// handle an update

	// catch all
	rw.WriteHeader(http.StatusNotImplemented)
}

func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle Get Products")
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

func (p *Products) addProduct(rw http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle POST Products")

	// about request ioreader,  it buffers data, go does not read all the content at once
	// ioreader reads data from http request chunk by chunk

	prod := &data.Product{}
	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
	}

	p.logger.Printf("Prod: %#v", prod)
	data.AddProduct(prod)

}
