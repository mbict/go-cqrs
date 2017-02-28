package dsl

import (
	"github.com/mbict/gogen/dslengine"
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen"
)

func Projection(name string, dsl func()) {
	domain, ok := domainDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	_, name = sanitzeMapKeyName(name)
	for _, p := range domain.Projections {
		if p.Name == name {
			dslengine.ReportError("duplicate projection declared, already one defined with the same name `%s`", name)
			return
		}
	}

	projection := &cqrs.ProjectionExpr{
		Name:   name,
		Domain: domain,
	}
	domain.Projections = append(domain.Projections, projection)
	dslengine.Execute(dsl, projection)
}

func HandlesEvents(event ...string) {
	projection, ok := projectionDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	em := map[string]bool{}
	emKeys := []string{}
	for _, eventName := range append(projection.HandlesEvent, event...) {
		if _, ok := em[eventName]; !ok {
			em[eventName] = true
			emKeys = append(emKeys, eventName)
		}
	}
	projection.HandlesEvent = emKeys
}

func Repository(name string, model *gogen.UserTypeExpr, dsl ... func()) {
	projection, ok := projectionDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}

	if len(dsl) > 1 {
		dslengine.ReportError("cannot have more than one dsl body, only the first is used")
	}

	_, name = sanitzeMapKeyName(name)
	for _, r := range projection.Domain.AllReadRepositories() {
		if r.Name == name {
			dslengine.ReportError("duplicate read repository declared, already one defined with the same name `%s`", name)
			return
		}
	}

	repo := &cqrs.RepositoryExpr{
		Name:  name,
		Model: model,
		Filter: &gogen.AttributeExpr{
			Type: &gogen.Object{},
		},
	}

	if len(dsl) >= 1 {
		dslengine.Execute(dsl[0], repo)
	}

	projection.ReadRepositories = append(projection.ReadRepositories, repo)
}

func Filter(dsl func()) {
	repo, ok := repositoryDefinition()
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}
	dslengine.Execute(dsl, repo.Attribute())
}
