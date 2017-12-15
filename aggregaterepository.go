package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
	"reflect"
)

// PublishEventFunc is a callback function that is getting called once the eventstore has successfully stored the
// new events gereated by the aggregate
type PublishEventFunc func(event Event)

// AggregateBuilder is the builder function to create new aggregate compositions.
// This could be used to introduce new strategies how to build a aggregate like the snapshot implementation
type AggregateBuilder func(aggregateId uuid.UUID) (AggregateComposition, error)

// AggregateRepository is the interface that a specific aggregate repositories should implement.
type AggregateRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateId uuid.UUID) (AggregateComposition, error)

	//Saves the aggregate.
	Save(aggregate AggregateComposition) error
}

type aggregateRepository struct {
	eventStore           EventStore
	aggregateBuilder     AggregateBuilder
	abstractEventFactory EventFactory
	publishEventHooks    []PublishEventFunc
}

func DefaultAggregateBuilder(factory AggregateFactoryFunc) AggregateBuilder {
	return func(aggregateId uuid.UUID) (AggregateComposition, error) {
		context := NewAggregateContext(aggregateId, 0)
		aggregateComposition := &aggregateContextComposition{
			AggregateContext: context,
		}
		aggregate := factory(aggregateComposition)
		if aggregate == nil {
			return nil, nil
		}

		aggregateComposition.Aggregate = aggregate
		return aggregateComposition, nil
	}
}

func (r *aggregateRepository) Load(aggregateId uuid.UUID) (AggregateComposition, error) {
	aggregate, err := r.aggregateBuilder(aggregateId)
	if err != nil {
		return nil, err
	}
	stream, err := r.eventStore.LoadStream(aggregate.AggregateName(), aggregateId, aggregate.Version())
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
		aggregate.incrementVersion()
	}
	return aggregate, nil
}

func (r *aggregateRepository) Save(aggregate AggregateComposition) error {
	events := aggregate.getUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	if err := r.eventStore.WriteEvent(aggregate.AggregateName(), events...); err != nil {
		return err
	}
	aggregate.clearUncommittedEvents()

	for _, event := range events {
		event = reflect.Indirect(reflect.ValueOf(event)).Interface().(Event)
		for _, f := range r.publishEventHooks {
			f(event)
		}
	}

	return nil
}

// NewAggregateRepository is the constructor of the repository
//
// publishEventHooks get called when a new event is successfully persisted to the eventstore.
// This is very useful to wire it to an eventbus for publishing the event to other listeners (projections)
func NewAggregateRepository(
	eventStore EventStore,
	aggregateBuilder AggregateBuilder,
	abstractEventFactory EventFactory,
	publishEventHooks ...PublishEventFunc) AggregateRepository {
	return &aggregateRepository{
		eventStore:           eventStore,
		aggregateBuilder:     aggregateBuilder,
		abstractEventFactory: abstractEventFactory,
		publishEventHooks:    publishEventHooks,
	}
}
