package memory

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
)

var ErrNoEventData = errors.New("Cannot scan, no event data")

type EventStream struct {
	yield        chan cqrs.Event
	currentEvent cqrs.Event
}

func newMemoryEventStream(events []cqrs.Event) cqrs.EventStream {
	yield := make(chan cqrs.Event)
	go func() {
		for _, event := range events {
			yield <- event
		}
		close(yield)
	}()

	return &EventStream{
		yield: yield,
	}
}

func (s *EventStream) EventName() string {
	if s.currentEvent == nil {
		return ""
	}
	return s.currentEvent.EventName()
}

func (s *EventStream) AggregateId() uuid.UUID {
	if s.currentEvent == nil {
		return uuid.Nil
	}
	return s.currentEvent.AggregateID()
}

func (s *EventStream) Version() int {
	if s.currentEvent == nil {
		return -1
	}
	return s.currentEvent.Version()
}

func (s *EventStream) Next() bool {
	s.currentEvent = <-s.yield
	return s.currentEvent != nil
}

func (s *EventStream) Error() error {
	return nil
}

func (s *EventStream) Scan(event cqrs.Event) error {
	if s.currentEvent == nil {
		return ErrNoEventData
	}
	return copier.Copy(event, s.currentEvent)
}
