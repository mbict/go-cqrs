package cqrs

import "github.com/mbict/gogen"

type CommandExpr struct {
	// Name of the command
	Name string

	//Params this command needs
	Params *gogen.AttributeExpr

	//Events are names of the events this commands generates
	Events []string

	// RootAggregate is the aggregate this command belongs to
	RootAggregate *AggregateExpr
}

func (c *CommandExpr) Context() string {
	return "command"
}
