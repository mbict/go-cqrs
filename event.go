package cqrs

import (
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

type Event interface {
	eventbus.Event
	AggregateID() uuid.UUID
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

func (e *EventBase) AggregateID() uuid.UUID {
	return e.id
}

func (e *EventBase) Version() int {
	return e.version
}
