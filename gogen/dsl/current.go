package dsl

import (
	. "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen/dslengine"
)

func domainDefinition() (*DomainExpr, bool) {
	d, ok := dslengine.Current().(*DomainExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}

func aggregateDefinition() (*AggregateExpr, bool) {
	a, ok := dslengine.Current().(*AggregateExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

func commandDefinition() (*CommandExpr, bool) {
	c, ok := dslengine.Current().(*CommandExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return c, ok
}

func projectionDefinition() (*ProjectionExpr, bool) {
	p, ok := dslengine.Current().(*ProjectionExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return p, ok
}

func repositoryDefinition() (*RepositoryExpr, bool) {
	p, ok := dslengine.Current().(*RepositoryExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return p, ok
}