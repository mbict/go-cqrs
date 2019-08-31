package cqrs

import (
	"time"
)

type EventStore interface {
	LoadStream(aggregateName string, aggregateId AggregateId, version int) (EventStream, error)
	WriteEvent(string, ...Event) error
}

type EventStream interface {
	EventType() EventType
	AggregateId() AggregateId
	Version() int
	Timestamp() time.Time

	Next() bool
	Error() error
	Scan(EventData) error
}
