package cqrs

import (
	"context"
	"fmt"
	"github.com/mbict/go-commandbus"
)

type ErrorAggregateNotFound string

func (e ErrorAggregateNotFound) Error() string {
	return fmt.Sprintf("aggregate not found for id: \"%s\"", string(e))
}

type ErrorNotAnAggregateCommand string

func (e ErrorNotAnAggregateCommand) Error() string {
	return fmt.Sprintf("cannot convert to aggregate command for command: \"%s\"", string(e))
}

type ErrorAggregateCannotHandleCommand string

func (e ErrorAggregateCannotHandleCommand) Error() string {
	return fmt.Sprintf("aggregate cannot handle commands directly for command: \"%s\"", string(e))
}

// AggregateCommandHandler is a command handler middleware who loads the aggregate
// calls the aggregate command handler to execute the business logic and saves the
// events to the aggregate store afterwards.
func AggregateCommandHandler(repository AggregateRepository) commandbus.CommandHandler {
	return AggregateCommandHandlerCallback(repository, func(aggregate Aggregate, cmd Command) error {
		agg, ok := aggregate.(AggregateHandlesCommands)
		if !ok {
			return ErrorAggregateCannotHandleCommand(cmd.CommandName())
		}
		return agg.HandleCommand(cmd)
	})
}

type AggregateCommandHandlerFunc func(aggregate Aggregate, command Command) error

func AggregateCommandHandlerCallback(repository AggregateRepository, handler AggregateCommandHandlerFunc) commandbus.CommandHandler {
	return commandbus.CommandHandlerFunc(func(_ context.Context, command commandbus.Command) error {
		cmd, ok := command.(Command)
		if !ok {
			return ErrorNotAnAggregateCommand(command.CommandName())
		}

		aggregate, err := repository.Load(cmd.AggregateId())
		if err != nil {
			return err
		}

		if aggregate == nil {
			return ErrorAggregateNotFound(cmd.AggregateId().String())
		}

		//run validation if there is a validate structure implemented
		if validate, ok := command.(Validate); ok {
			if err := validate.Validate(); err != nil {
				return err
			}
		}

		if err = handler(aggregate, cmd); err != nil {
			return err
		}

		if err = repository.Save(aggregate); err != nil {
			return err
		}
		return nil
	})
}
