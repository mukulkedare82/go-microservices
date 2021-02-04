package handlers

import (
	"log"
	"net/http"
	"regexp"
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

func (p *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		p.GetProducts(rw, req)
		return
	}

	if req.Method == http.MethodPost {
		p.AddProduct(rw, req)
		return
	}

	if req.Method == http.MethodPut {
		p.logger.Println("Handling Put Products")
		// expect the id in the URI
		path := req.URL.Path
		regex := regexp.MustCompile("/([0-9]+)")
		group := regex.FindAllStringSubmatch(path, -1)

		if len(group) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			p.logger.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.logger.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
			return
		}

		p.logger.Println("got id", id)
		p.updateProduct(id, rw, req)
		return
	}

	// handle an update

	// catch all
	rw.WriteHeader(http.StatusNotImplemented)
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
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "unable to convert id", http.StatusBadRequest)
	}

	p.updateProduct(id, rw, req) // kept it as wrapper to work with servemux handler ServeHTTP

	return
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(req.Body)
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
