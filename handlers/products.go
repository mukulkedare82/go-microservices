package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mukulkedare/go-microservice-tuts/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) GetProducts(rw http.ResponseWriter, req *http.Request) {
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

func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
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

func (p *Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle PUT Product")

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "unable to convert id", http.StatusBadRequest)
	}

	prod := &data.Product{}
	err = prod.FromJSON(req.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
	}

	p.logger.Printf("Prod: %#v", prod)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product update failed", http.StatusInternalServerError)
		return
	}

	return

}
