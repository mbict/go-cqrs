package dsl

import (
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen/dslengine"
)

func Domain(name string, dsl func()) *cqrs.DomainExpr {
	_, name = sanitzeMapKeyName(name)
	for _, d := range cqrs.Root {
		if d.Name == name {
			dslengine.ReportError("duplicate domain declared, already one defined with the same name `%s`", name)
			return nil
		}
	}

	d := &cqrs.DomainExpr{
		Name: name,
	}
	dslengine.Execute(dsl, d)

	cqrs.Root = append(cqrs.Root, d)
	return d
}
