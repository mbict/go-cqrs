package cqrs

import (
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

// Event is the interface of an event what an aggregate needs
type Event interface {
	eventbus.Event
	AggregateId() uuid.UUID
	Version() int
}

// EventBase is an utility class for not reimplementing AggregateId and Version
// methods of the Event interface
type EventBase struct {
	id      uuid.UUID
	version int
}

// NewEventBase constructor with plain version
func NewEventBase(id uuid.UUID, version int) *EventBase {
	return &EventBase{
		id:      id,
		version: version,
	}
}

// NewEventBaseFromAggregate constructor will create a new eventbase
// based on the latest aggregate state
func NewEventBaseFromAggregate(aggregate Aggregate) *EventBase {
	return &EventBase{
		id:      aggregate.AggregateId(),
		version: aggregate.CurrentVersion() + 1,
	}
}

// AggregateId returns the id of the aggregate
func (e EventBase) AggregateId() uuid.UUID {
	return e.id
}

// Version returns the event version/sequence in the stream
func (e EventBase) Version() int {
	return e.version
}
