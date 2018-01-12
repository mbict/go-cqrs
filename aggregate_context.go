package cqrs

import (
	"github.com/satori/go.uuid"
)

type AggregateContext interface {
	// AggregateId returns the id of the aggregate.
	AggregateId() uuid.UUID

	// Version returns the version of the aggregate.
	Version() int

	// StoreEvent stores an event as uncommitted event.
	StoreEvent(Event)

	// incrementVersion increments the aggregate version.
	incrementVersion()

	// setVersion is an internal function to set the aggregate version.
	setVersion(int)

	// getUncommittedEvents gets all uncommitted events ready for storing.
	getUncommittedEvents() []Event

	// clearUncommittedEvents clears all uncommitted events after storing.
	clearUncommittedEvents()
}

type aggregateContext struct {
	id      uuid.UUID
	version int
	events  []Event
}

func NewAggregateContext(id uuid.UUID, version int) AggregateContext {
	return &aggregateContext{
		id:      id,
		events:  []Event{},
		version: version,
	}
}

func (a *aggregateContext) AggregateId() uuid.UUID {
	return a.id
}

func (a *aggregateContext) Version() int {
	return a.version
}

func (a *aggregateContext) StoreEvent(event Event) {
	a.events = append(a.events, event)
}

func (a *aggregateContext) incrementVersion() {
	a.version++
}

func (a *aggregateContext) setVersion(version int) {
	a.version = version
}

func (a *aggregateContext) getUncommittedEvents() []Event {
	return a.events
}

func (a *aggregateContext) clearUncommittedEvents() {
	a.events = []Event{}
}
