package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// AggregateBuilder is the builder function to create new aggregate compositions.
// This could be used to introduce new strategies how to build a aggregate like the snapshot implementation
type AggregateBuilder func(aggregateId uuid.UUID) (Aggregate, error)

// AggregateRepository is the interface that a specific aggregate repositories should implement.
type AggregateRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateId uuid.UUID) (Aggregate, error)

	//Saves the aggregate.
	Save(aggregate Aggregate) error
}

type aggregateRepository struct {
	eventStore       EventStore
	aggregateBuilder AggregateBuilder
	eventFactory     EventFactory
}

func DefaultAggregateBuilder(factory AggregateFactoryFunc) AggregateBuilder {
	return func(aggregateId uuid.UUID) (Aggregate, error) {
		context := NewAggregateContext(aggregateId, 0)
		aggregate := factory(context)
		if aggregate == nil {
			return nil, nil
		}
		return aggregate, nil
	}
}

func (r *aggregateRepository) Load(aggregateId uuid.UUID) (Aggregate, error) {
	aggregate, err := r.aggregateBuilder(aggregateId)
	if err != nil {
		return nil, err
	}
	stream, err := r.eventStore.LoadStream(aggregate.AggregateName(), aggregateId, aggregate.Version())
	if err != nil {
		return nil, fmt.Errorf("cannot load events from stream reader, error: %s", err)
	}

	for stream != nil && stream.Next() {
		if stream.Version()-1 != aggregate.Version() {
			return nil, fmt.Errorf("event version (%d) mismatch with Aggregate next Version (%d)", stream.Version(), aggregate.Version()+1)
		}

		eventData := r.eventFactory.MakeEvent(stream.EventType())
		if eventData == nil {
			return nil, fmt.Errorf("the repository has no event factory registered for event type: %s", stream.EventType())
		}

		err = stream.Scan(eventData)
		if err != nil {
			return nil, fmt.Errorf("the repository cannot populate event data from stream for event type: %s, with error `%s`", stream.EventType(), err)
		}

		//create the event with metadata
		event := NewEvent(aggregateId, stream.Version(), stream.Timestamp(), eventData)

		aggregate.Apply(event)
		aggregate.incrementVersion()
	}
	return aggregate, nil
}

func (r *aggregateRepository) Save(aggregate Aggregate) error {
	events := aggregate.getUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	if err := r.eventStore.WriteEvent(aggregate.AggregateName(), events...); err != nil {
		return err
	}
	aggregate.clearUncommittedEvents()

	//apply events
	for _, event := range events {
		aggregate.Apply(event)
		aggregate.incrementVersion()
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
	eventFactory EventFactory) AggregateRepository {
	return &aggregateRepository{
		eventStore:       eventStore,
		aggregateBuilder: aggregateBuilder,
		eventFactory:     eventFactory,
	}
}
