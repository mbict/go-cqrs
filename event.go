package cqrs

import "github.com/satori/go.uuid"

type Event interface {
	AggregateID() uuid.UUID
	EventType() string
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
