package cqrs

import "github.com/satori/go.uuid"

type EventStore interface {
	LoadStream(aggregateType string, aggregateId uuid.UUID) (EventStream, error)
	FindStream(aggregateTypes []string, aggregateIds []uuid.UUID, eventTypes []string) (EventStream, error)
	WriteEvent(string, Event) error
}

type EventStoreAggregate interface {
}

type EventStream interface {
	EventType() string
	AggregateId() uuid.UUID
	Version() int

	Next() bool
	Error() error
	Scan(Event) error
}
