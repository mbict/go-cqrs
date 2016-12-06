package cqrs

type CommandHandler interface {
	HandleCommand(Command) error
}

// CommandBus is an interface defining an event bus for distributing events.
type CommandBus interface {
	// HandleCommand handles a command on the event bus.
	HandleCommand(Command) error
	// SetHandler registers a handler with a command.
	SetHandler(CommandHandler, Command) error
}
