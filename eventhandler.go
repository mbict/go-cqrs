package cqrs

import "github.com/mbict/go-eventbus"

type EventHandler interface {
	Handle(Event) error
}

type EventHandlerFunc func(event Event) error

func (h EventHandlerFunc) Handle(event Event) error {
	return h(event)
}

func EventbusWrapper(handler EventHandler) eventbus.EventHandlerFunc {
	return func(event eventbus.Event) error {
		if e, ok := event.(Event); ok {
			return handler.Handle(e)
		}
		return nil
	}
}
