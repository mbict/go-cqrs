package cqrs

import "github.com/mbict/gogen"

type EventExpr struct {
	// Name of the event
	Name string

	// Attributes this event has
	Attributes *gogen.AttributeExpr

	// RootAggregate where this event is generated
	RootAggregate *AggregateExpr
}

func (e *EventExpr) Context() string {
	return "event"
}


func (e *EventExpr) Attribute() *gogen.AttributeExpr {
	return e.Attributes
}
