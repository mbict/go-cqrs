package cqrs

import (
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

type Event interface {
	eventbus.Event
	AggregateId() uuid.UUID
	Version() int
}

type EventBase struct {
	id      uuid.UUID
	version int
}

func NewEventBase(id uuid.UUID, version int) *EventBase {
	return &EventBase{
		id:      id,
		version: version,
	}
}

func NewEventBaseFromAggregate(aggregate Aggregate) *EventBase {
	return &EventBase{
		id:      aggregate.AggregateId(),
		version: aggregate.CurrentVersion() + 1,
	}
}

func (e EventBase) AggregateId() uuid.UUID {
	return e.id
}

func (e EventBase) Version() int {
	return e.version
}
