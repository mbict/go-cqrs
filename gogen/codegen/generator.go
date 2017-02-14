package codegen

import (
	"errors"
	"fmt"
	"path"
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen"
	"text/template"
)

type Generator struct {
	gogen.CodeGenerator

	basePackage string
}

func NewGenerator(basePackage string) *Generator {
	cg := gogen.NewCodeGenerator("./templates/*.tmpl")

	return &Generator{
		CodeGenerator: cg,
		basePackage:   basePackage,
	}
}

func (g *Generator) Writers(root interface{}) ([]gogen.FileWriter, error) {
	domain, ok := root.(*cqrs.DomainExpr)
	if !ok {
		return nil, fmt.Errorf("Incompatible root")
	}

	res := []gogen.FileWriter{}

	for _, a := range domain.Aggregates {
		s, err := g.GenerateAggregate(a)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("domain/aggregate/%s.go", a.Name)
		fwProjection := gogen.NewFileWriter(s, file)
		res = append(res, fwProjection)

		for _, c := range a.Commands {
			s, err := g.GenerateCommand(c)
			if err != nil {
				return nil, err
			}

			file := fmt.Sprintf("domain/commands/%s.go", c.Name)
			fwCommand := gogen.NewFileWriter(s, file)
			res = append(res, fwCommand)
		}
	}

	for _, e := range domain.AllEvents() {
		s, err := g.GenerateEvent(e)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("domain/events/%s.go", e.Name)
		fwEvent := gogen.NewFileWriter(s, file)
		res = append(res, fwEvent)
	}

	for _, p := range domain.Projections {
		s, err := g.GenerateProjection(p)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("domain/projection/%s.go", p.Name)
		fwProjection := gogen.NewFileWriter(s, file)
		res = append(res, fwProjection)
	}

	return res, nil
}

func (g *Generator) GenerateEvent(e *cqrs.EventExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("EVENT")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(path.Join(g.basePackage, "events"))
	imports.Add("github.com/mbict/go-cqrs")
	imports.AddFromAttribute(e.Attributes)

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Event":   e,
			"Imports": imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateCommand(c *cqrs.CommandExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("COMMAND")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports("")
	imports.Add(gogen.UUID.Package)
	imports.AddFromAttribute(c.Params)

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Command": c,
			"Imports": imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateAggregate(a *cqrs.AggregateExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("AGGREGATE")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(path.Join(g.basePackage, "Aggregates"))
	imports.Add("errors")
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "events"))
	imports.Add(path.Join(g.basePackage, "commands"))

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Aggregate": a,
			"Imports":   imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateProjection(p *cqrs.ProjectionExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("PROJECTION")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(path.Join(g.basePackage, "Projections"))
	imports.Add("errors")
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(path.Join(g.basePackage, "events"))

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Projection": p,
			"Imports":    imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateAggregatesFactory(d *cqrs.DomainExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("AGGREGATE_FACTORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(g.basePackage)
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "aggregates"))

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateEventsFactory(d *cqrs.DomainExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("EVENT_FACTORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(g.basePackage)
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "events"))

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateRepositoryInterfaces(d *cqrs.DomainExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("REPOSITORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(g.basePackage)
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "models"))
	for _, r := range d.AllReadRepositories() {
		if r.Filter != nil {
			imports.AddFromAttribute(r.Filter)
		}
	}

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}

	return []gogen.Section{s}, nil
}

func (g *Generator) GenerateDbRepository(r *cqrs.RepositoryExpr) ([]gogen.Section, error) {
	t := g.Template().Lookup("DB_REPOSITORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := gogen.NewImports(g.basePackage)
	imports.Add("database/sql")
	imports.Add(gogen.UUID.Package)
	imports.Add("github.com/masterminds/squirrel")
	imports.Add(path.Join(g.basePackage, "models"))
	if r.Filter != nil {
		imports.AddFromAttribute(r.Filter)
	}

	s := gogen.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Repository":  r,
			"Imports":     imports,
		},
	}

	return []gogen.Section{s}, nil
}
