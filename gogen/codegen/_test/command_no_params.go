package command

import (
	"github.com/satori/go.uuid"
)

type Test struct {
	Id uuid.UUID
}

func (c *Test) AggregateId() uuid.UUID {
	return c.Id
}

func (c *Test) AggregateType() string {
	return "Item"
}

func (c *Test) CommandType() string {
	return "Test"
}
