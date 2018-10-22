package cqrs

import (
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
	"time"
)

// Event is the interface of an event what an aggregate needs
type Event interface {
	eventbus.Event
	EventBase
}

type EventBase interface {
	AggregateId() uuid.UUID
	Version() int
	OccurredAt() time.Time
}

// EventBase is an utility class for not reimplementing AggregateId and Version
// methods of the Event interface
type eventBase struct {
	id         uuid.UUID
	version    int
	occurredAt time.Time
}

// NewEventBase constructor with plain version
func NewEventBase(id uuid.UUID, version int, occurredAt time.Time) EventBase {
	return &eventBase{
		id:         id,
		version:    version,
		occurredAt: occurredAt,
	}
}

// NewEventBaseFromAggregate constructor will create a new eventbase
// based on the latest aggregate state
func NewEventBaseFromAggregate(aggregate AggregateContext) EventBase {
	return &eventBase{
		id:         aggregate.AggregateId(),
		version:    aggregate.OriginalVersion() + len(aggregate.getUncommittedEvents()) + 1,
		occurredAt: time.Now(),
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

// OccurredAt returns the date and time the event occurred the first time
func (e *eventBase) OccurredAt() time.Time {
	return e.occurredAt
}
