package memory

import (
	"fmt"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
)

type EventStore struct {
	events map[string]map[string][]cqrs.Event
}

func NewMemoryEventStore() cqrs.EventStore {
	return &EventStore{
		events: make(map[string]map[string][]cqrs.Event),
	}
}

func (s *EventStore) LoadStream(aggregateType string, aggregateId uuid.UUID) (cqrs.EventStream, error) {
	aggregates, ok := s.events[aggregateType]
	if !ok {
		return nil, nil
	}

	events, ok := aggregates[aggregateId]
	if !ok {
		return nil, nil
	}
	return newMemoryEventStream(events), nil
}

func (s *EventStore) WriteEvent(aggregateType string, event cqrs.Event) error {

	if _, ok := s.events[aggregateType]; !ok {
		s.events[aggregateType] = make(map[string][]cqrs.Event)
	}
	s.events[aggregateType][event.AggregateID()] = append(s.events[aggregateType][event.AggregateID()], event)
	fmt.Printf("Saving event %s for aggregate %s (%s)\n", event.EventType(), aggregateType, event.AggregateID())
	return nil
}
