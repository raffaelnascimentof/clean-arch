package events

import "time"

// EventInterface defines the contract for domain events.
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

// EventHandlerInterface defines the contract for event handlers.
type EventHandlerInterface interface {
	Handle(event EventInterface, wg interface{})
}

// EventDispatcherInterface defines the contract for the event dispatcher.
type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear() error
}
