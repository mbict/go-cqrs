package cqrs

import "github.com/satori/go.uuid"

type EventStore interface {
	LoadStream(aggregateName string, aggregateId uuid.UUID) (EventStream, error)
	WriteEvent(string, ...Event) error
}

type EventStream interface {
	EventName() string
	AggregateId() uuid.UUID
	Version() int

	Next() bool
	Error() error
	Scan(Event) error
}
