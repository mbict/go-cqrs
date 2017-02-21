package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// DomainRepository is the interface that all domain repositories should implement.
type DomainRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateTypeName string, aggregateId uuid.UUID) (AggregateRoot, error)

	//Saves the aggregate.
	Save(aggregate AggregateRoot) error
}

type CommonDomainRepository struct {
	eventStore       EventStore
	eventFactory     EventFactory
	aggregateFactory AggregateFactory
	eventBus         EventBus
}

//NewRepository instantiates a new repository resolver who accepts a stream resolver
func NewCommonDomainRepository(eventStore EventStore, eventFactory EventFactory, aggregateFactory AggregateFactory) *CommonDomainRepository {
	return &CommonDomainRepository{
		eventBus:         nil,
		eventStore:       eventStore,
		eventFactory:     eventFactory,
		aggregateFactory: aggregateFactory,
	}
}

//SetEventBus will set wich eventbus for publishing new events
//if set to nil no events will be published during a save action
func (r *CommonDomainRepository) SetEventBus(eventBus EventBus) {
	r.eventBus = eventBus
}

//Loads an aggregate of the given type and ID
func (r *CommonDomainRepository) Load(aggregateType string, aggregateId uuid.UUID) (AggregateRoot, error) {

	aggregate := r.aggregateFactory.GetAggregate(aggregateType, aggregateId)
	if aggregate == nil {
		return nil, fmt.Errorf("The repository has no aggregate factory registered for aggregate type: %s", aggregateType)
	}

	stream, err := r.eventStore.LoadStream(aggregateType, aggregateId)
	if err != nil {
		return nil, fmt.Errorf("Cannot load events from stream reader for aggregate type: %s, error: %s", aggregateType, err)
	}

	for stream != nil && stream.Next() {

		event := r.eventFactory.GetEvent(stream.EventType(), aggregateId, stream.Version())
		if event == nil {
			return nil, fmt.Errorf("The repository has no event factory registered for event type: %s", stream.EventType())
		}

		if stream.Version()-1 != aggregate.Version() {
			return nil, fmt.Errorf("Event version (%d) mismatch with Aggregate next Version (%d)", stream.Version(), aggregate.Version()+1)
		}

		err = stream.Scan(event)
		if err != nil {
			return nil, fmt.Errorf("The repository cannot populate event data from stream for event type: %s", stream.EventType())
		}

		//aggregate.Apply(event)
		//aggregate.IncrementVersion()
	}

	//snapshotversion = 0
	//if snapshotversion-aggrevgate.Version() > 50 {
	//	//create snapshot
	//}


	return aggregate, nil
}

//Save will save all the events to the event store.
func (r *CommonDomainRepository) Save(aggregate AggregateRoot) error {

	for _, event := range aggregate.GetUncommittedEvents() {
		if err := r.eventStore.WriteEvent(aggregate.AggregateType(), event); err != nil {
			return err
		}

		if r.eventBus != nil {
			r.eventBus.PublishEvent(event)
		}

		aggregate.Apply(event)

	}

	aggregate.ClearUncommittedEvents()

	//snapshot strategy
	return nil
}
