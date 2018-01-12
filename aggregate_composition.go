package cqrs

type AggregateComposition interface {
	Context() AggregateContext
	Aggregate() Aggregate
}
