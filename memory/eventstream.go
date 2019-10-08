package memory

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/mbict/go-cqrs/v4"
	"time"
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

func (s *EventStream) EventType() cqrs.EventType {
	if s.currentEvent == nil {
		return ""
	}
	return s.currentEvent.EventType()
}

func (s *EventStream) AggregateId() cqrs.AggregateId {
	if s.currentEvent == nil {
		return uuid.Nil
	}
	return s.currentEvent.AggregateId()
}

func (s *EventStream) Version() int {
	if s.currentEvent == nil {
		return -1
	}
	return s.currentEvent.Version()
}

var emptyTime = time.Time{}

func (s *EventStream) Timestamp() time.Time {
	if s.currentEvent == nil {
		return emptyTime
	}
	return s.currentEvent.Timestamp()
}

func (s *EventStream) Next() bool {
	s.currentEvent = <-s.yield
	return s.currentEvent != nil
}

func (s *EventStream) Error() error {
	return nil
}

func (s *EventStream) Scan(event cqrs.EventData) error {
	if s.currentEvent == nil {
		return ErrNoEventData
	}
	return copier.Copy(event, s.currentEvent.Data())
}
