package cqrs

import (
	"github.com/satori/go.uuid"
)

// aggregateSnapshotContextComposition is a wrapper to store the initial snapshot version
type aggregateSnapshotContextComposition struct {
	snapshotVersion int
	AggregateContext
	Aggregate
}

func SnapshotAggregateBuilder(factory AggregateFactoryFunc, snapshotStore SnapshotStore) AggregateBuilder {
	return func(aggregateId uuid.UUID) (AggregateComposition, error) {
		aggregateComposition := &aggregateSnapshotContextComposition{}
		aggregateState := factory(aggregateComposition)

		version, err := snapshotStore.Load(aggregateId, aggregateState)
		if err != nil {
			return nil, err
		}

		aggregateComposition.snapshotVersion = version
		aggregateComposition.Aggregate = aggregateState
		aggregateComposition.AggregateContext = NewAggregateContext(aggregateId, version)

		return aggregateComposition, nil
	}
}

type snapshotAggregateRepository struct {
	AggregateRepository
	snapshotStore    SnapshotStore
	differenceOffset int
}

func (r *snapshotAggregateRepository) Save(aggregate AggregateComposition) error {
	if err := r.AggregateRepository.Save(aggregate); err != nil {
		return err
	}

	//if this aggregate is constructed out of a snapshot composition we check if we need to create a new snapshot
	if aggSnapshotComp, ok := aggregate.(*aggregateSnapshotContextComposition); ok {
		needSnapshot := (aggSnapshotComp.snapshotVersion + r.differenceOffset) < aggregate.Version()
		if needSnapshot == true {
			return r.snapshotStore.Write(aggregate)
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
	abstractEventFactory EventFactory,
	publishEventHooks ...PublishEventFunc) AggregateRepository {
	return &snapshotAggregateRepository{
		snapshotStore:    snapshotStore,
		differenceOffset: differenceOffset,
		AggregateRepository: &aggregateRepository{
			eventStore:           eventStore,
			aggregateBuilder:     aggregateBuilder,
			abstractEventFactory: abstractEventFactory,
			publishEventHooks:    publishEventHooks,
		},
	}
}
