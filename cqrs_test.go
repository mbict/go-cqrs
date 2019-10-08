package cqrs

import (
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
	Id AggregateId
}

func (*commandA) CommandName() string {
	return "commandA"
}

func (c *commandA) AggregateId() AggregateId {
	return c.Id
}

type eventA struct {
	Id AggregateId
}

func (e eventA) AggregateId() AggregateId {
	return e.Id
}

func (e eventA) Version() int {
	panic("implement me")
}

func (e eventA) Timestamp() time.Time {
	panic("implement me")
}

func (e eventA) Data() EventData {
	panic("implement me")
}

func (eventA) EventType() EventType {
	return "event:a"
}

type eventB struct {
	Id AggregateId
}

func (e eventB) AggregateId() AggregateId {
	return e.Id
}

func (e eventB) Version() int {
	panic("implement me")
}

func (e eventB) Timestamp() time.Time {
	panic("implement me")
}

func (e eventB) Data() EventData {
	panic("implement me")
}

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
