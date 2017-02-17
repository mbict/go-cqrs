package dsl

import (
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/go-cqrs/gogen/codegen"
	"github.com/mbict/gogen/dslengine"
)

func Domain(name string, dsl func()) *cqrs.DomainExpr {
	if cqrs.Root != nil {
		dslengine.ReportError("duplicate domain declared, can have only one", name)
		return nil
	}

	codegen.Register() //register codegenerator

	_, name = sanitzeMapKeyName(name)
	d := &cqrs.DomainExpr{
		Name: name,
		Events: make(map[string]*cqrs.EventExpr, 0),
	}
	dslengine.Execute(dsl, d)

	cqrs.Root = d
	return d
}
