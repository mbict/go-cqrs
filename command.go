package cqrs

import (
	"github.com/satori/go.uuid"
)

type Command interface {
	CommandName() string
	AggregateId() uuid.UUID
}

type AggregateCommand interface {
	Command
	AggregateName() string
}
