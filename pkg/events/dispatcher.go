package events

import (
	"errors"
	"sync"
)

// ErrHandlerAlreadyRegistered is returned when a handler is already registered for an event.
var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

// EventDispatcher manages event registrations and dispatching.
type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

// NewEventDispatcher creates a new EventDispatcher instance.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

// Register adds a handler for a given event name.
func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

// Dispatch fires all registered handlers for the given event concurrently.
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

// Remove unregisters a handler for a given event name.
func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

// Has checks whether a handler is registered for a given event name.
func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

// Clear removes all registered handlers.
func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}
