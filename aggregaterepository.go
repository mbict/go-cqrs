package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
	"reflect"
)

// AggregateRepository is the interface that a specific aggregate repositories should implement.
type AggregateRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateId uuid.UUID) (Aggregate, error)

	//Saves the aggregate.
	Save(aggregate Aggregate) error
}

type PublishEventFunc func(event Event)

type aggregateRepository struct {
	eventStore           EventStore
	aggregateFactory     AggregateFactoryFunc
	abstractEventFactory EventFactory
	publishEventHooks    []PublishEventFunc
}

func (r *aggregateRepository) Load(aggregateId uuid.UUID) (Aggregate, error) {
	aggregate := r.aggregateFactory(aggregateId)
	stream, err := r.eventStore.LoadStream(aggregate.AggregateName(), aggregateId)
	if err != nil {
		return nil, fmt.Errorf("cannot load events from stream reader, error: %s", err)
	}

	for stream != nil && stream.Next() {
		event := r.abstractEventFactory.MakeEvent(stream.EventName(), aggregateId, stream.Version())
		if event == nil {
			return nil, fmt.Errorf("the repository has no event factory registered for event type: %s", stream.EventName())
		}

		if stream.Version()-1 != aggregate.Version() {
			return nil, fmt.Errorf("event version (%d) mismatch with Aggregate next Version (%d)", stream.Version(), aggregate.Version()+1)
		}

		err = stream.Scan(event)
		if err != nil {
			return nil, fmt.Errorf("the repository cannot populate event data from stream for event type: %s, with error `%s`", stream.EventName(), err)
		}

		// we do not want to pass events by pointer reference but by pass by value,
		// just to ensure the data of the events are readonly so no other process can change them
		event = reflect.Indirect(reflect.ValueOf(event)).Interface().(Event)
		aggregate.Apply(event)
		aggregate.IncrementVersion()
	}

	return aggregate, nil
}

func (r *aggregateRepository) Save(aggregate Aggregate) error {
	events := aggregate.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	if err := r.eventStore.WriteEvent(aggregate.AggregateName(), events...); err != nil {
		return err
	}
	aggregate.ClearUncommittedEvents()

	for _, event := range events {
		event = reflect.Indirect(reflect.ValueOf(event)).Interface().(Event)
		aggregate.Apply(event)
		aggregate.IncrementVersion()

		for _, f := range r.publishEventHooks {
			f(event)
		}
	}
	return nil
}

func NewAggregateRepository(
	eventStore EventStore,
	aggregateFactory AggregateFactoryFunc,
	abstractEventFactory EventFactory,
	publishEventHooks ...PublishEventFunc) AggregateRepository {
	return &aggregateRepository{
		eventStore:           eventStore,
		aggregateFactory:     aggregateFactory,
		abstractEventFactory: abstractEventFactory,
		publishEventHooks:    publishEventHooks,
	}
}
