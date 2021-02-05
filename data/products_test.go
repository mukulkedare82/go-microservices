package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "mukul",
		Price: 1.00,
		//SKU:   "abc",
		SKU: "abc-absd-dfsdf",
	}

	err := p.Validate()

	if err != nil {

		t.Fatal(err)
	}
}
