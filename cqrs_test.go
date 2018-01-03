package cqrs

import (
	"github.com/satori/go.uuid"
	"time"
)

type commandAWithValidate struct {
	commandA
	ValidateMock *MockValidate
}

func (m *commandAWithValidate) Validate() error {
	return m.ValidateMock.Validate()
}

type commandNonAggregate struct {
}

func (*commandNonAggregate) CommandName() string {
	return "commandNonAggregate"
}

type commandA struct {
	Id uuid.UUID
}

func (*commandA) CommandName() string {
	return "commandA"
}

func (c *commandA) AggregateId() uuid.UUID {
	return c.Id
}

type eventA struct {
	EventBase
}

func (eventA) EventName() string {
	return "event.a"
}

type eventB struct {
	EventBase
}

func (eventB) EventName() string {
	return "event.b"
}

func eventAFactory(id uuid.UUID, version int, occurredAt time.Time) Event {
	return &eventA{
		EventBase: NewEventBase(id, version, occurredAt),
	}
}

type aggregateA struct {
	ctx AggregateContext
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

func aggregateAFactory(ctx AggregateContext) Aggregate {
	return &aggregateA{
		ctx: ctx,
	}
}
