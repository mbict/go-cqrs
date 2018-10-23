package cqrs

import (
	"github.com/satori/go.uuid"
	"time"
)

type EventStore interface {
	LoadStream(aggregateName string, aggregateId uuid.UUID, version int) (EventStream, error)
	WriteEvent(string, ...Event) error
}

type EventStream interface {
	EventType() EventType
	AggregateId() uuid.UUID
	Version() int
	Timestamp() time.Time

	Next() bool
	Error() error
	Scan(EventData) error
}
