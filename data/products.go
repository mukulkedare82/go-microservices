package data

import (
	"encoding/json"
	"io"
	"time"
)

// Product defines structure payload for API
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// about "struct tags" in go
// allows adding annotations to struct fields
// during parsing we can use these annotations
// Ex:
// json is the struct tag for json parsing. we can name, omitempty, ignore, as key

// about "slice of struct" in go
// we can define a type which is slice of another struct
// kind of derived/custom type
type Products []*Product

// json encode Products (slice of products)
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p) // returns an Error, writes data directly to io.Wrtier
}

func GetProducts() Products {
	return productList
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}