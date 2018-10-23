package cqrs

type publishingEventStore struct {
	publisher EventPublisher
	EventStore
}

func (s *publishingEventStore) WriteEvent(name string, events ...Event) error {
	err := s.EventStore.WriteEvent(name, events...)
	if err != nil {
		return err
	}
	return s.publisher.Publish(events...)
}

func NewEventPublishingEventStore(publisher EventPublisher, store EventStore) EventStore {
	return &publishingEventStore{
		publisher:  publisher,
		EventStore: store,
	}
}
