package codegen

import (
	"errors"
	"fmt"
	"path"
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/codegen/generator"
	"text/template"
	"runtime"
	"path/filepath"
)

type Generator struct {
}

type Codegen struct {
	gogen.CodeGenerator
	basePackage string
}

var gen *Generator

func init() {
	gen = &Generator{}
}

func Register() {
	generator.Register(gen)
}


func (g *Generator) Name() string {
	return "cqrs"
}

func (g *Generator) Generate(path string) error {
	//todo : remove go path
	codegen := NewCodeGenerator(path)
	fileWriters, err := codegen.Writers( cqrs.Root )
	if err != nil {
		return err
	}

	for _, fw := range fileWriters {
		contents, err := fw.Write()
		if err != nil {
			return err
		}
		fmt.Println("==================")
		fmt.Println(fw.Path())
		fmt.Println("==================")
		fmt.Println(contents)
	}
	return nil
}


func NewCodeGenerator(basePackage string) *Codegen {
	_, file, _, _ := runtime.Caller(0)
	templatePath := filepath.Join(path.Dir(file), "templates", "*.tmpl")
	cg := gogen.NewCodeGenerator(templatePath)
	return &Codegen{
		CodeGenerator: cg,
		basePackage: basePackage,
	}
}

func (g *Codegen) Writers(root interface{}) ([]gogen.FileWriter, error) {
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

func (g *Codegen) GenerateEvent(e *cqrs.EventExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateCommand(c *cqrs.CommandExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateAggregate(a *cqrs.AggregateExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateProjection(p *cqrs.ProjectionExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateAggregatesFactory(d *cqrs.DomainExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateEventsFactory(d *cqrs.DomainExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateRepositoryInterfaces(d *cqrs.DomainExpr) ([]gogen.Section, error) {
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

func (g *Codegen) GenerateDbRepository(r *cqrs.RepositoryExpr) ([]gogen.Section, error) {
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
