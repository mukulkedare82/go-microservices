package handlers

import (
	"net/http"
	"regexp"
	"strconv"
)

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
