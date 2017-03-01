package event

import (
	"github.com/mbict/go-cqrs"
)

type Testevent1 struct {
	*cqrs.EventBase
}

func (e *Testevent1) AggregateType() string {
	return "Item"
}

func (e *Testevent1) EventType() string {
	return "Testevent1"
}
