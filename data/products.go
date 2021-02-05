package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product defines structure payload for API
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
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

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// create a validator instance
func (p *Product) Validate() error {
	validate := validator.New()

	// register custom validators
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

// create a custom validator for sku field
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-dfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)

	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true

}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)

}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
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
