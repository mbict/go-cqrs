package cqrs

import (
	"fmt"
)

type ErrorAggregateFactoryAlreadyRegistered string

func (e ErrorAggregateFactoryAlreadyRegistered) Error() string {
	return fmt.Sprintf("aggregate factory callback/delegate already registered for type: \"%s\"", string(e))
}

type AggregateFactoryFunc func(AggregateContext) Aggregate

// AggregateFactory returns aggregate instances of a specified type with the AggregateId set to the uuid provided.
type AggregateFactory interface {
	//MakeAggregate will return a clean Aggregate based on the type provided
	MakeAggregate(string, AggregateContext) Aggregate
}

// CallbackAggregateFactory is an implementation of the AggregateFactory interface
// that supports registration of delegate/callback functions to perform aggregate instantiation.
type CallbackAggregateFactory struct {
	delegates map[string]AggregateFactoryFunc
}

// NewCallbackAggregateFactory creates a new CallbackAggregateFactory
func NewCallbackAggregateFactory() *CallbackAggregateFactory {
	return &CallbackAggregateFactory{
		delegates: make(map[string]AggregateFactoryFunc),
	}
}

// RegisterCallback is used to register a new function for instantiation of an aggregate instance.
func (t *CallbackAggregateFactory) RegisterCallback(callback AggregateFactoryFunc) error {
	typeName := callback(nil).AggregateName()
	if _, ok := t.delegates[typeName]; ok {
		return ErrorAggregateFactoryAlreadyRegistered(typeName)
	}
	t.delegates[typeName] = callback
	return nil
}

// MakeAggregate calls the callback for the specified type and returns the result.
func (t *CallbackAggregateFactory) MakeAggregate(typeName string, ctx AggregateContext) Aggregate {
	if f, ok := t.delegates[typeName]; ok {
		return f(ctx)
	}
	return nil
}
