package cqrs

import (
	"github.com/satori/go.uuid"
)

type eventA struct {
	*EventBase
}

func (*eventA) EventName() string {
	return "event.a"
}

func eventAFactory(id uuid.UUID, version int) Event {
	return &eventA{
		EventBase: NewEventBase(id, version),
	}
}

type aggregateA struct {
	*AggregateBase
}

func (*aggregateA) AggregateName() string {
	return "aggregateA"
}

func (*aggregateA) HandleCommand(Command) error {
	return nil
}

func (*aggregateA) Apply(event Event) error {
	return nil
}

func aggregateAFactory(id uuid.UUID) Aggregate {
	return &aggregateA{
		AggregateBase: NewAggregateBase(id),
	}
}
