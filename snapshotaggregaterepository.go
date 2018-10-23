package cqrs

import (
	"github.com/satori/go.uuid"
)

// aggregateSnapshotComposition is a wrapper to store the initial snapshot version
type aggregateSnapshotComposition struct {
	snapshotVersion int
	aggregate       Aggregate
}

func (c *aggregateSnapshotComposition) Context() AggregateContext {
	return c.aggregate
}

func (c *aggregateSnapshotComposition) Aggregate() Aggregate {
	if a, ok := c.aggregate.(AggregateComposition); ok {
		return a.Aggregate()
	}
	return c.aggregate
}

func (c *aggregateSnapshotComposition) AggregateId() uuid.UUID {
	return c.aggregate.AggregateId()
}

func (c *aggregateSnapshotComposition) Version() int {
	return c.aggregate.Version()
}

func (c *aggregateSnapshotComposition) OriginalVersion() int {
	return c.aggregate.OriginalVersion()
}

func (c *aggregateSnapshotComposition) StoreEvent(e EventData) {
	c.aggregate.StoreEvent(e)
}

func (c *aggregateSnapshotComposition) incrementVersion() {
	c.aggregate.incrementVersion()
}

func (c *aggregateSnapshotComposition) setVersion(version int) {
	c.aggregate.setVersion(version)
}

func (c *aggregateSnapshotComposition) getUncommittedEvents() []Event {
	return c.aggregate.getUncommittedEvents()
}

func (c *aggregateSnapshotComposition) clearUncommittedEvents() {
	c.aggregate.clearUncommittedEvents()
}

func (c *aggregateSnapshotComposition) AggregateName() string {
	return c.aggregate.AggregateName()
}

func (c *aggregateSnapshotComposition) Apply(e Event) error {
	return c.aggregate.Apply(e)
}

// SnapshotAggregateBuilder
func SnapshotAggregateBuilder(factory AggregateFactoryFunc, snapshotStore SnapshotStore) AggregateBuilder {
	return func(aggregateId uuid.UUID) (Aggregate, error) {
		aggregateComposition := &aggregateSnapshotComposition{}
		context := NewAggregateContext(aggregateId, 0)
		aggregate := factory(context)

		version, err := snapshotStore.Load(aggregateId, aggregate)
		if err != nil {
			return nil, err
		}

		aggregate.setVersion(version)
		aggregateComposition.snapshotVersion = version
		aggregateComposition.aggregate = aggregate

		return aggregateComposition, nil
	}
}

type snapshotAggregateRepository struct {
	AggregateRepository
	snapshotStore    SnapshotStore
	differenceOffset int
}

func (r *snapshotAggregateRepository) Save(aggregate Aggregate) error {
	if err := r.AggregateRepository.Save(aggregate); err != nil {
		return err
	}

	//if this aggregate is constructed out of a snapshot composition we check if we need to create a new snapshot
	if aggSnapshotComp, ok := aggregate.(*aggregateSnapshotComposition); ok {
		needSnapshot := aggregate.Version() >= (aggSnapshotComp.snapshotVersion + r.differenceOffset)
		if needSnapshot == true {
			return r.snapshotStore.Write(aggSnapshotComp.Aggregate())
		}
	}

	return nil
}

// NewSnapshotAggregateRepository is the constructor of the aggregate repository with snapshot functionality
// A snapshot will be created when the differenceOffset between the snapshot version and the current version is equal
// or greater than the `differenceOffset`
//
// When the differenceOffset is set to 10 than:
// - aggregate version 7 (snapshot version 0) will not create a snapshot
// - aggregate version 10 (snapshot version 0) will create a snapshot for version 10
// - aggregate version 13 (snapshot version 0) will create a snapshot for version 13
// - aggregate version 21 (snapshot version 13) will not create a snapshot
// - aggregate version 54 (snapshot version 13) will create a snapshot for version 54
func NewSnapshotAggregateRepository(
	eventStore EventStore,
	snapshotStore SnapshotStore,
	differenceOffset int,
	aggregateBuilder AggregateBuilder,
	eventFactory EventFactory) AggregateRepository {
	return &snapshotAggregateRepository{
		snapshotStore:    snapshotStore,
		differenceOffset: differenceOffset,
		AggregateRepository: &aggregateRepository{
			eventStore:       eventStore,
			aggregateBuilder: aggregateBuilder,
			eventFactory:     eventFactory,
		},
	}
}
