package cqrs

import "github.com/satori/go.uuid"

type SnapshotStore interface {
	Load(aggregateId uuid.UUID, aggregate Aggregate) (int, error)
	Write(aggregate AggregateComposition) error
}
