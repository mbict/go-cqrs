package cqrs

type Command interface {
	CommandName() string
	AggregateId() AggregateId
}

type AggregateCommand interface {
	Command
	AggregateName() string
}
