package cqrs

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// AggregateFactory returns aggregate instances of a specified type with the AggregateID set to the uuid provided.
type AggregateFactory interface {
	//GetAggregate will return a clean AggregateRoot based on the type provided
	GetAggregate(string, uuid.UUID) AggregateRoot
}

// CallbackAggregateFactory is an implementation of the AggregateFactory interface
// that supports registration of delegate/callback functions to perform aggregate instantiation.
type CallbackAggregateFactory struct {
	delegates map[string]func(uuid.UUID) AggregateRoot
}

// NewCallbackAggregateFactory creates a new CallbackAggregateFactory
func NewCallbackAggregateFactory() *CallbackAggregateFactory {
	return &CallbackAggregateFactory{
		delegates: make(map[string]func(uuid.UUID) AggregateRoot),
	}
}

// RegisterCallback is used to register a new function for instantiation of an aggregate instance.
func (t *CallbackAggregateFactory) RegisterCallback(callback func(uuid.UUID) AggregateRoot) error {
	typeName := callback(uuid.NewV4()).AggregateType()
	if _, ok := t.delegates[typeName]; ok {
		return fmt.Errorf("Factory callback/delegate already registered for type: \"%s\"", typeName)
	}
	t.delegates[typeName] = callback
	return nil
}

// GetAggregate calls the callback for the specified type and returns the result.
func (t *CallbackAggregateFactory) GetAggregate(typeName string, id uuid.UUID) AggregateRoot {
	if f, ok := t.delegates[typeName]; ok {
		return f(id)
	}
	return nil
}
