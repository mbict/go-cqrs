package cqrs

import (
	"errors"
	"fmt"
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

// AggregateRepositoryManager is the managing interface who provide aggregate repository access
type AggregateRepositoryManager interface {
	//RepositoryFor will return the repository for the specific named aggregate
	RepositoryFor(aggregateName string) AggregateRepository
}

type DomainAggregateRepository struct {
	eventStore       EventStore
	eventFactory     EventFactory
	aggregateFactory AggregateFactory
	eventBus         eventbus.EventBus
}

//NewRepository instantiates a new repository resolver who accepts a stream resolver
func NewCommonDomainRepository(eventStore EventStore, eventFactory EventFactory, aggregateFactory AggregateFactory) *DomainAggregateRepository {
	return &DomainAggregateRepository{
		eventBus:         nil,
		eventStore:       eventStore,
		eventFactory:     eventFactory,
		aggregateFactory: aggregateFactory,
	}
}

//SetEventBus will set which eventbus for publishing new events
//if set to nil no events will be published during a save action
func (r *DomainAggregateRepository) SetEventBus(eventBus eventbus.EventBus) {
	r.eventBus = eventBus
}

func (r *DomainAggregateRepository) RepositoryFor(aggregateName string) AggregateRepository {
	panic("implement me")
}

//Loads an aggregate of the given type and ID
func (r *DomainAggregateRepository) Load(aggregateType string, aggregateId uuid.UUID) (Aggregate, error) {
	context := NewAggregateContext(aggregateId, 0)
	aggregate := newAggregateContextComposition(context, func(ctx AggregateContext) Aggregate {
		return r.aggregateFactory.MakeAggregate(aggregateType, context)
	})

	if aggregate == nil {
		return nil, fmt.Errorf("the repository has no aggregate factory registered for aggregate type: %s", aggregateType)
	}

	stream, err := r.eventStore.LoadStream(aggregateType, aggregateId, context.Version())
	if err != nil {
		return nil, fmt.Errorf("cannot load events from stream reader for aggregate type: %s, error: %s", aggregateType, err)
	}

	for stream != nil && stream.Next() {

		event := r.eventFactory.MakeEvent(stream.EventName(), aggregateId, stream.Version(), stream.OccurredAt())
		if event == nil {
			return nil, fmt.Errorf("the repository has no event factory registered for event type: %s", stream.EventName())
		}

		if stream.Version()-1 != aggregate.Version() {
			return nil, fmt.Errorf("event version (%d) mismatch with Aggregate next Version (%d)", stream.Version(), aggregate.Version()+1)
		}

		err = stream.Scan(event)
		if err != nil {
			return nil, fmt.Errorf("the repository cannot populate event data from stream for event type: %s", stream.EventName())
		}

		aggregate.Apply(event)
	}

	return aggregate, nil
}

//Save will save all the events to the event store.
func (r *DomainAggregateRepository) Save(aggregate AggregateComposition) error {
	for _, event := range aggregate.getUncommittedEvents() {
		if err := r.eventStore.WriteEvent(aggregate.AggregateName(), event); err != nil {
			return err
		}

		aggregate.Apply(event)

		if r.eventBus != nil {
			r.eventBus.Publish(event)
		}
	}
	aggregate.clearUncommittedEvents()
	return nil
}
