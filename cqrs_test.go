package cqrs

import (
	"github.com/satori/go.uuid"
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
}

func (eventA) EventType() EventType {
	return "event:a"
}

type eventB struct{}

func (eventB) EventType() EventType {
	return "event.b"
}

func eventAFactory() EventData {
	return &eventA{}
}

type aggregateA struct {
	AggregateContext
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
		AggregateContext: ctx,
	}
}
