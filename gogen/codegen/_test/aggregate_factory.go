package domain

import (
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"testing/base/aggregates"
)

func RegisterAggregateFactories(factory *cqrs.CallbackAggregateFactory) error {
	testAggregate1AggregateFactory := func(aggregateId uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewTestAggregate1Aggregate(aggregateId)
	}
	if err := factory.RegisterCallback(testAggregate1AggregateFactory); err != nil {
		return err
	}

	testAggregate2AggregateFactory := func(aggregateId uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewTestAggregate2Aggregate(aggregateId)
	}
	if err := factory.RegisterCallback(testAggregate2AggregateFactory); err != nil {
		return err
	}

	return nil
}
