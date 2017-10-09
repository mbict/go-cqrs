package cqrs

import (
	"github.com/mbict/go-commandbus"
	"github.com/satori/go.uuid"
)

type Command interface {
	commandbus.Command
	AggregateId() uuid.UUID
}

type AggregateCommand interface {
	Command
	AggregateName() string
}
