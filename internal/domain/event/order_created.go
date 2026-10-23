package event

import "time"

// OrderCreated is the domain event fired when an order is successfully created.
type OrderCreated struct {
	Name    string
	Payload interface{}
}

// NewOrderCreated instantiates an OrderCreated event.
func NewOrderCreated() *OrderCreated {
	return &OrderCreated{Name: "OrderCreated"}
}

// GetName returns the event name.
func (o *OrderCreated) GetName() string { return o.Name }

// GetDateTime returns the current time as the event timestamp.
func (o *OrderCreated) GetDateTime() time.Time { return time.Now() }

// GetPayload returns the event payload.
func (o *OrderCreated) GetPayload() interface{} { return o.Payload }

// SetPayload sets the event payload.
func (o *OrderCreated) SetPayload(payload interface{}) { o.Payload = payload }
