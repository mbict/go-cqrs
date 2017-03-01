package event

import (
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
)

type Testevent1 struct {
	*cqrs.EventBase
	Test1 string
	Test2 []uuid.UUID
}

func (e *Testevent1) AggregateType() string {
	return "Item"
}

func (e *Testevent1) EventType() string {
	return "Testevent1"
}
