// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
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

	// get value from req context by key
	prod := req.Context().Value(KeyProduct{}).(*data.Product)

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

	// about fetching params from context
	// get value from req context by key
	// returns interface, pass type to it for casting return value
	prod := req.Context().Value(KeyProduct{}).(*data.Product)

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

// about using gorilla middleware
// Use function allows adding middleware (MiddlewareFunc) to a route
// middleware is nothing but http handler
// with middleware pattern we can chain multiple handlers together
// Ex: applying middleware to handle CORS validation

type KeyProduct struct{}

func (p *Products) MiddleWareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		p.logger.Println("Inside MiddleWareProductionValidation")
		// validate request json and create prod object
		prod := &data.Product{}
		err := prod.FromJSON(req.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserialising product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product using the struct tagged validator
		err = prod.Validate()
		if err != nil {
			p.logger.Println("[ERROR] validating product", err)
			http.Error(rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest)
			return
		}

		// create request context and pass product object for next handler in chain
		// create key(type struct or string) for the prod and pass it as key/val pair in context
		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req2 := req.WithContext(ctx)

		next.ServeHTTP(rw, req2)
	})
}
