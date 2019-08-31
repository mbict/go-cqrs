package cqrs

type SnapshotStore interface {
	Load(aggregateId AggregateId, aggregate Aggregate) (int, error)
	Write(aggregate Aggregate) error
}
