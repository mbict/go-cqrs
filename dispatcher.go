package cqrs

import (
	"errors"
)

var (
	ErrNilRepository     = errors.New("Repository cannot be nil")
	ErrAggregateNotFound = errors.New("Aggregate not found")
)

type Dispatcher struct {
	repository DomainRepository
}

// NewDispatcher creates a new Dispatcher who will route commands to the right aggregate command handler
func NewDispatcher(repository DomainRepository) (*Dispatcher, error) {
	if repository == nil {
		return nil, ErrNilRepository
	}

	h := &Dispatcher{
		repository: repository,
	}
	return h, nil
}

// HandleCommand handles a command with the registered aggregate.
// Returns ErrAggregateNotFound if no aggregate could be found.
func (h *Dispatcher) HandleCommand(command Command) error {

	aggregate, err := h.repository.Load(command.AggregateType(), command.AggregateId())
	if err != nil {
		return err
	}

	if aggregate == nil {
		return ErrAggregateNotFound
	}

	//run validation if there is a validate structure implemented
	if validate, ok	:= command.(Validate); ok {
		if err := validate.Validate(); err != nil {
			return err
		}
	}

	if err = aggregate.HandleCommand(command); err != nil {
		return err
	}

	if err = h.repository.Save(aggregate); err != nil {
		return err
	}

	return nil
}
