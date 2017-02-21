package dsl

import (
	"github.com/mbict/gogen/dslengine"
	"github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen"
)

func Aggregate(name string, dsl func()) {
	domain, ok := domainDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	_, name = sanitzeMapKeyName(name)
	for _, a := range domain.Aggregates {
		if a.Name == name {
			dslengine.ReportError("duplicate aggregate declared, already one defined with the same name `%s`", name)
			return
		}
	}

	aggregate := &cqrs.AggregateExpr{
		Name:   name,
		Domain: domain,
	}
	domain.Aggregates = append(domain.Aggregates, aggregate)
	dslengine.Execute(dsl, aggregate)
}

func Command(name string, dsl ...func()) {
	aggregate, ok := aggregateDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	if len(dsl) > 1 {
		dslengine.ReportError("cannot have more than one dsl body, only the first is used")
	}

	_, name = sanitzeMapKeyName(name)
	for _, c := range aggregate.Commands {
		if c.Name == name {
			dslengine.ReportError("duplicate command declared, already one defined with the same name `%s`", name)
			return
		}
	}

	command := &cqrs.CommandExpr{
		Name: name,
		Params: &gogen.AttributeExpr{
			Type: &gogen.Object{
				&gogen.Field{
					Name: "Id",
					Attribute: &gogen.AttributeExpr{
						Type: gogen.UUID,
						Validation: &gogen.ValidationRule{
							Required: true,
						},
					},
				},
			},
		},
		RootAggregate: aggregate,
	}

	if len(dsl) >= 1 {
		dslengine.Execute(dsl[0], command)
	}
	aggregate.Commands = append(aggregate.Commands, command)
}

func Params(dsl func()) {
	command, ok := commandDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	dslengine.Execute(dsl, command.Params)
}

func Event(name string, dsl ...func()) {
	command, ok := commandDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	if len(dsl) > 1 {
		dslengine.ReportError("cannot have more than one dsl body, only the first is used")
	}

	_, name = sanitzeMapKeyName(name)
	aggregate := command.RootAggregate
	domain := aggregate.Domain
	if _, ok := domain.Events[name]; ok {
		dslengine.ReportError("duplicate event declared, already one defined with the same name `%s`", name)
	}

	event := &cqrs.EventExpr{
		Name: name,
		Attributes: &gogen.AttributeExpr{
			Type: &gogen.Object{
				&gogen.Field{
					Name: "Id",
					Attribute: &gogen.AttributeExpr{
						Type: gogen.UUID,
						Validation: &gogen.ValidationRule{
							Required: true,
						},
					},
				},
			},
		},
		RootAggregate: aggregate,
	}

	command.Events = append(command.Events, name)

	if len(dsl) >= 1 {
		dslengine.Execute(dsl[0], event)
	}
	domain.Events[name] = event
}
