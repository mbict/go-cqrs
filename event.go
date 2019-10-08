package cqrs

import (
	"github.com/mbict/go-eventbus"
	"reflect"
	"time"
)

type EventType = eventbus.EventType

// Event
type Event interface {
	EventData
	Version() int
	Timestamp() time.Time
	Data() interface{}
}

// EventData is the actual data of the event
type EventData interface {
	AggregateId() AggregateId
	EventType() EventType
}

// Event is an utility class for not reimplementing AggregateId and Version
// methods of the Event interface
type event struct {
	EventData
	version   int
	timestamp time.Time
}

func (e *event) Timestamp() time.Time {
	return e.timestamp
}

// Timestamp returns the date and time the event occurred the first time
func (e *event) Data() interface{} {
	return e.EventData
}

// Version returns the event version/sequence in the stream
func (e *event) Version() int {
	return e.version
}

// NewEvent constructor with plain version
func NewEvent(version int, timestamp time.Time, data EventData) Event {
	return &event{
		EventData: passEventByValue(data),
		version:   version,
		timestamp: timestamp,
	}
}

// NewEventFromAggregate constructor will create a new event
// based on the latest aggregate state
func NewEventFromAggregate(aggregate AggregateContext, data EventData) Event {
	return &event{
		EventData: passEventByValue(data),
		version:   aggregate.OriginalVersion() + len(aggregate.getUncommittedEvents()) + 1,
		timestamp: time.Now(),
	}
}

// we do not want to pass events by pointer reference but by pass by value,
// just to ensure the data of the events are readonly so no other process can change them
func passEventByValue(data EventData) EventData {
	return reflect.Indirect(reflect.ValueOf(data)).Interface().(EventData)
}
