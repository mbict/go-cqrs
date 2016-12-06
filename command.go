package cqrs

import "github.com/satori/go.uuid"

type Command interface {
	AggregateId() uuid.UUID
	AggregateType() string
	CommandType() string
}
