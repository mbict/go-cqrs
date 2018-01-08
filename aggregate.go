package cqrs

type Aggregate interface {
	// AggregateName returns the type name of the aggregate.
	AggregateName() string

	// Apply applies an event to the aggregate by setting its values.
	Apply(Event) error
}

//AggregateHandlesCommands indicates a aggregate can directly handle a command
type AggregateHandlesCommands interface {
	Aggregate

	// HandleCommand handles a command and stores events.
	HandleCommand(Command) error
}
