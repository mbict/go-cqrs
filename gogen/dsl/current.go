package dsl

import (
	. "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen/dslengine"
)

func domainDefinition() (*DomainExpr, bool) {
	d, ok := dslengine.CurrentDefinition().(*DomainExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}

func aggregateDefinition() (*AggregateExpr, bool) {
	a, ok := dslengine.CurrentDefinition().(*AggregateExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return a, ok
}

func commandDefinition() (*CommandExpr, bool) {
	c, ok := dslengine.CurrentDefinition().(*CommandExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return c, ok
}

func projectionDefinition() (*ProjectionExpr, bool) {
	p, ok := dslengine.CurrentDefinition().(*ProjectionExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return p, ok
}

func repositoryDefinition() (*RepositoryExpr, bool) {
	p, ok := dslengine.CurrentDefinition().(*RepositoryExpr)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return p, ok
}