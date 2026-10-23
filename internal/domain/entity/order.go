package entity

import (
	"errors"

	"github.com/google/uuid"
)

// Order represents the Order entity in the domain layer.
type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

// NewOrder creates a new Order, generating a UUID and calculating the final price.
func NewOrder(id string, price float64, tax float64) (*Order, error) {
	if id == "" {
		id = uuid.NewString()
	}
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	if err := order.Validate(); err != nil {
		return nil, err
	}
	order.CalculateFinalPrice()
	return order, nil
}

// Validate performs basic validation on the Order fields.
func (o *Order) Validate() error {
	if o.Price <= 0 {
		return errors.New("invalid price: must be greater than zero")
	}
	if o.Tax < 0 {
		return errors.New("invalid tax: must be non-negative")
	}
	return nil
}

// CalculateFinalPrice computes the final price including tax.
func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
