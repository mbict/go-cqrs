package memory

import (
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"sync"
)

type EventStore struct {
	events  map[string]map[string][]cqrs.Event
	rwMutex sync.RWMutex
}

func NewMemoryEventStore() cqrs.EventStore {
	return &EventStore{
		events: make(map[string]map[string][]cqrs.Event),
	}
}

func (s *EventStore) LoadStream(aggregateName string, aggregateId uuid.UUID, version int) (cqrs.EventStream, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	aggregates, ok := s.events[aggregateName]
	if !ok {
		return nil, nil
	}

	events, ok := aggregates[aggregateId.String()]
	if !ok {
		return nil, nil
	}

	return newMemoryEventStream(events[version:]), nil
}

func (s *EventStore) WriteEvent(aggregateName string, events ...cqrs.Event) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for _, event := range events {
		if _, ok := s.events[aggregateName]; !ok {
			s.events[aggregateName] = make(map[string][]cqrs.Event)
		}
		id := event.AggregateId().String()
		s.events[aggregateName][id] = append(s.events[aggregateName][id], event)
	}
	return nil
}
