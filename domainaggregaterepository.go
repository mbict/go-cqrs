package cqrs

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

// AggregateRepositoryManager is the managing interface who provide aggregate repository access
type AggregateRepositoryManager interface {
	//RepositoryFor will return the repository for the specific named aggregate
	RepositoryFor(aggregateName string) (AggregateRepository, error)

	//Loads an aggregate based on the aggregate ID and aggregateName
	Load(aggregateName string, aggregateId AggregateId) (Aggregate, error)

	//Saves the aggregate.
	Save(aggregate Aggregate) error
}

type DomainAggregateRepository struct {
	eventStore            EventStore
	eventFactory          EventFactory
	aggregateFactory      AggregateFactory
	aggregateRepositories map[string]AggregateRepository
}

//NewRepository instantiates a new repository resolver who accepts a stream resolver
func NewCommonDomainRepository(eventStore EventStore, eventFactory EventFactory, aggregateFactory AggregateFactory) *DomainAggregateRepository {
	return &DomainAggregateRepository{
		eventStore:            eventStore,
		eventFactory:          eventFactory,
		aggregateFactory:      aggregateFactory,
		aggregateRepositories: make(map[string]AggregateRepository),
	}
}

func (r *DomainAggregateRepository) RepositoryFor(aggregateName string) (AggregateRepository, error) {
	if repository, ok := r.aggregateRepositories[aggregateName]; ok {
		return repository, nil
	}

	aggregateBuilder := r.aggregateBuilderFor(aggregateName)

	//try to build the aggregate as validation
	if agg, err := aggregateBuilder(uuid.Nil); err != nil || agg == nil {
		return nil, ErrRepositoryNotFound
	}

	//cache the aggregateRepo and return the instance
	r.aggregateRepositories[aggregateName] = NewAggregateRepository(r.eventStore, aggregateBuilder, r.eventFactory)
	return r.aggregateRepositories[aggregateName], nil
}

func (r *DomainAggregateRepository) aggregateBuilderFor(aggregateName string) AggregateBuilder {
	return func(aggregateId AggregateId) (Aggregate, error) {
		context := NewAggregateContext(aggregateId, 0)
		return r.aggregateFactory.MakeAggregate(aggregateName, context), nil
	}
}

//Loads an aggregate of the given type and ID
func (r *DomainAggregateRepository) Load(aggregateName string, aggregateId AggregateId) (Aggregate, error) {
	repository, err := r.RepositoryFor(aggregateName)
	if err != nil {
		return nil, err
	}
	return repository.Load(aggregateId)
}

//Save will save all the events to the event store.
func (r *DomainAggregateRepository) Save(aggregate Aggregate) error {
	repository, err := r.RepositoryFor(aggregate.AggregateName())
	if err != nil {
		return err
	}
	return repository.Save(aggregate)
}
