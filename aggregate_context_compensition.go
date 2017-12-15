package cqrs

import (
	"encoding/json"
)

type AggregateComposition interface {
	AggregateContext
	Aggregate
}

// aggregateContextComposition is a wrapper object that ensures the apply function of the
// aggregate is triggered upon a StoreEvent call
type aggregateContextComposition struct {
	AggregateContext
	Aggregate
}

func (c *aggregateContextComposition) StoreEvent(event Event) {
	c.AggregateContext.StoreEvent(event)
	c.AggregateContext.incrementVersion()
	c.Aggregate.Apply(event)
}

func (c *aggregateContextComposition) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Aggregate)
}

func (c *aggregateContextComposition) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.Aggregate)
}

func newAggregateContextComposition(context AggregateContext, factory AggregateFactoryFunc) AggregateComposition {
	aggregateComposition := &aggregateContextComposition{
		AggregateContext: context,
	}
	aggregate := factory(aggregateComposition)
	if aggregate == nil {
		return nil
	}

	aggregateComposition.Aggregate = aggregate
	return aggregateComposition
}
