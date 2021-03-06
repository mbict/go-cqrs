package cqrs

import (
	"testing"
)

func TestCallbackAggregateFactoryMake(t *testing.T) {
	f := NewCallbackAggregateFactory()

	if err := f.RegisterCallback(aggregateAFactory); err != nil {
		t.Errorf("expected a nil error but got error : %v", err)
	}

	id := NewStringAggregateId("abc")
	ctx := NewAggregateContext(id, 0)

	aggregate := f.MakeAggregate("aggregateA", ctx)
	if aggregate == nil {
		t.Fatal("expected the constructed aggregate but got nil instead")
	}

	if aggregate.AggregateName() != "aggregateA" {
		t.Errorf("expected an aggregate with name `%v` but got `%v`", "aggregateA", aggregate.AggregateName())
	}

	if aggregate.AggregateId().String() != "abc" {
		t.Errorf("expected an aggregate with id `%v` but got `%v`", "abc", aggregate.AggregateId())
	}

	if aggregate.Version() != 0 {
		t.Errorf("expected an aggregate with version `%v` but got `%v`", 0, aggregate.Version())
	}

	if aggregate.AggregateName() != "aggregateA" {
		t.Errorf("expected an aggregate with name `%v` but got `%v`", "aggregateA", aggregate.AggregateName())
	}

	if _, ok := aggregate.(*aggregateA); !ok {
		t.Errorf("expected an aggregate of type `%v` but got `%T`", "aggregateA", aggregate)
	}
}

func TestCallbackAggregateFactoryMakeWithUnknownAggregate(t *testing.T) {
	f := NewCallbackAggregateFactory()

	ctx := NewAggregateContext(NewStringAggregateId("abc"), 0)
	aggregate := f.MakeAggregate("this.aggregate.is.not.registered", ctx)
	if aggregate != nil {
		t.Fatalf("expected a nil response but got an aggregate instead `%T`", aggregate)
	}
}

func TestCallbackAggregateFactoryDuplicateRegister(t *testing.T) {
	f := NewCallbackAggregateFactory()

	if err := f.RegisterCallback(aggregateAFactory); err != nil {
		t.Errorf("expected a nil error but got error : %v", err)
	}

	err := f.RegisterCallback(aggregateAFactory)
	if err == nil {
		t.Error("expected a error but got none")
	}

	if _, ok := err.(ErrorAggregateFactoryAlreadyRegistered); !ok {
		t.Error("expected a error but got none")
	}
}

func TestErrorAggregateFactoryAlreadyRegisteredToString(t *testing.T) {
	e := ErrorAggregateFactoryAlreadyRegistered("test")

	if e.Error() == "" {
		t.Errorf("expected a error message `%s` but got message : `%s`", "", e.Error())
	}
}
