package cqrs

type Aggregate interface {
	// AggregateName returns the type name of the aggregate.
	AggregateName() string

	// HandleCommand handles a command and stores events.
	HandleCommand(Command) error

	// Apply applies an event to the aggregate by setting its values.
	Apply(Event) error
}
