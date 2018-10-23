package cqrs

type EventPublisher interface {
	Publish(...Event) error
}

type EventPublisherFunc func(event ...Event) error

func (h EventPublisherFunc) Publish(event ...Event) error {
	return h(event...)
}
