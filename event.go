package cqrs

import (
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

// Event is the interface of an event what an aggregate needs
type Event interface {
	eventbus.Event
	EventBase
}

type EventBase interface {
	AggregateId() uuid.UUID
	Version() int
}

// EventBase is an utility class for not reimplementing AggregateId and Version
// methods of the Event interface
type eventBase struct {
	id      uuid.UUID
	version int
}

// NewEventBase constructor with plain version
func NewEventBase(id uuid.UUID, version int) EventBase {
	return &eventBase{
		id:      id,
		version: version,
	}
}

// NewEventBaseFromAggregate constructor will create a new eventbase
// based on the latest aggregate state
func NewEventBaseFromAggregate(aggregate AggregateContext) EventBase {
	return &eventBase{
		id:      aggregate.AggregateId(),
		version: aggregate.Version() + 1,
	}
}

// AggregateId returns the id of the aggregate
func (e *eventBase) AggregateId() uuid.UUID {
	return e.id
}

// Version returns the event version/sequence in the stream
func (e *eventBase) Version() int {
	return e.version
}
