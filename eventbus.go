package cqrs

// EventHandler is the interface that an subscriber to the event bus must implement.
type EventHandler interface {
	// HandleEvent handles an event.
	HandleEvent(Event)
}

// EventBus is the inteface that an event bus must implement.
type EventBus interface {
	PublishEvent(Event)
	AddHandler(EventHandler, ...Event)
}

type localEventBus struct {
	handlers map[string][]EventHandler
	catchAll []EventHandler
}

// NewLocalEventBus creates a net local in memory eventbus
func NewLocalEventBus() EventBus {
	return &localEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *localEventBus) PublishEvent(event Event) {
	if handlers, ok := b.handlers[event.EventType()]; ok {
		for _, handler := range handlers {
			handler.HandleEvent(event)
		}
	}
}

func (b *localEventBus) AddHandler(handler EventHandler, events ...Event) {
	for _, event := range events {
		if !b.hasHandlerForEvent(event.EventType(), handler) {
			b.handlers[event.EventType()] = append(b.handlers[event.EventType()], handler)
		}
	}
}

func (b *localEventBus) hasHandlerForEvent(event string, handler EventHandler) bool {

	handlers, ok := b.handlers[event]
	if !ok {
		return false
	}

	for _, h := range handlers {
		if h == handler {
			return true
		}
	}
	return false
}
