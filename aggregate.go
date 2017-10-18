package cqrs

import "github.com/satori/go.uuid"

type Aggregate interface {
	// AggregateId returns the id of the aggregate.
	AggregateId() uuid.UUID

	// AggregateName returns the type name of the aggregate.
	AggregateName() string

	// Version returns the version of the aggregate.
	Version() int

	// CurrentVersion returns the new version of the aggregate including the uncomitted events.
	CurrentVersion() int

	// IncrementVersion increments the aggregate version.
	IncrementVersion()

	// HandleCommand handles a command and stores events.
	HandleCommand(Command) error

	// Apply applies an event to the aggregate by setting its values.
	Apply(Event) error

	// StoreEvent stores an event as uncommitted event.
	StoreEvent(Event)

	// GetUncommittedEvents gets all uncommitted events ready for storing.
	GetUncommittedEvents() []Event

	// ClearUncommittedEvents clears all uncommitted events after storing.
	ClearUncommittedEvents()
}

// AggregateBase is a type that can be embedded in an Aggregate
// implementation to handle common aggregate behaviour
//
// All required methods to implement an aggregate are here, to implement the
// Aggregate root interface your aggregate will need to implement the Apply
// and the handle method that will contain behaviour specific logic
// to your aggregate.
type AggregateBase struct {
	id      uuid.UUID
	version int
	events  []Event
}

// NewAggregateBase constructs a new AggregateBase.
func NewAggregateBase(id uuid.UUID) *AggregateBase {
	return &AggregateBase{
		id:      id,
		events:  []Event{},
		version: 0,
	}
}

func (a *AggregateBase) AggregateId() uuid.UUID {
	return a.id
}

func (a *AggregateBase) Version() int {
	return a.version
}

func (a *AggregateBase) CurrentVersion() int {
	return a.version + len(a.events)
}

func (a *AggregateBase) IncrementVersion() {
	a.version++
}

func (a *AggregateBase) StoreEvent(event Event) {
	a.events = append(a.events, event)
}

func (a *AggregateBase) GetUncommittedEvents() []Event {
	return a.events
}

func (a *AggregateBase) ClearUncommittedEvents() {
	a.events = []Event{}
}
