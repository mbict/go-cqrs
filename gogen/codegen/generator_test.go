package codegen

import (
	. "gopkg.in/check.v1"
	cqrs "github.com/mbict/go-cqrs/gogen"
	"github.com/mbict/gogen"
	"testing"
	"bytes"
	"os"
	"io"
)

func Test(t *testing.T) { TestingT(t) }

type GeneratorSuite struct {
	generator *Codegen

	domain *cqrs.DomainExpr
}

var _ = Suite(&GeneratorSuite{})

func (s *GeneratorSuite) SetUpTest(c *C) {
	s.generator = NewCodeGenerator("testing/base")
	s.domain = &cqrs.DomainExpr{
		Name: "maindomain",
		Events: map[string]*cqrs.EventExpr{
			"ItemCreated": &cqrs.EventExpr{
				Name: "ItemCreated",
			},
			"ItemPriceUpdated": &cqrs.EventExpr{
				Name: "ItemPriceUpdated",
			},
			"ItemTitleUpdated": &cqrs.EventExpr{
				Name: "ItemTitleUpdated",
				Attributes: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "attr1",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
							},
						},
					},
				},
			},
			"ItemDeleted": &cqrs.EventExpr{
				Name: "ItemDeleted",
			},
		},
	}
}

func (s *GeneratorSuite) TestGenerateCommand(c *C) {

	rootAggregate := &cqrs.AggregateExpr{
		Name:   "item",
		Domain: s.domain,
	}

	tests := []struct {
		name     string
		command  *cqrs.CommandExpr
		expected string
	}{
		{
			name: "command with no params",
			command: &cqrs.CommandExpr{
				Name:   "test",
				Events: []string{},
				Params: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "id",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.UUID,
								Validation: &gogen.ValidationRule{
									Required: true,
								},
							},
						},
					},
				},
				RootAggregate: rootAggregate,
			},
			expected: loadfile(`_test/command_no_params.go`),
		}, {
			name: "command with params",
			command: &cqrs.CommandExpr{
				Name: "test",
				Params: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "id",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.UUID,
								Validation: &gogen.ValidationRule{
									Required: true,
								},
							},
						},
						&gogen.Field{
							Name: "name",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
								Validation: &gogen.ValidationRule{
									Required: true,
								},
							},
						},
						&gogen.Field{
							Name: "test",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
								Validation: &gogen.ValidationRule{
									Required: false,
								},
							},
						},
						&gogen.Field{
							Name: "tags",
							Attribute: &gogen.AttributeExpr{
								Type: &gogen.Array{
									ElemType: &gogen.AttributeExpr{
										Type: gogen.UUID,
									},
								},
								Validation: &gogen.ValidationRule{
									Required: false,
								},
							},

						},
					},
				},
				RootAggregate: rootAggregate,
			},
			expected: loadfile(`_test/command_params.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateCommand(test.command)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateEvent(c *C) {

	rootAggregate := &cqrs.AggregateExpr{
		Name:   "item",
		Domain: s.domain,
	}

	tests := []struct {
		name     string
		event    *cqrs.EventExpr
		expected string
	}{
		{
			name: "event with no attributes",
			event: &cqrs.EventExpr{
				Name:          "testevent1",
				RootAggregate: rootAggregate,
			},
			expected: loadfile(`_test/event_no_attributes.go`),
		}, {
			name: "event with params",
			event: &cqrs.EventExpr{
				Name: "testevent1",
				Attributes: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "test1",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
							},
						},
						&gogen.Field{
							Name: "test2",
							Attribute: &gogen.AttributeExpr{
								Type: &gogen.Array{
									ElemType: &gogen.AttributeExpr{
										Type: gogen.UUID,
									},
								},
							},
						},
					},
				},
				RootAggregate: rootAggregate,
			},
			expected: loadfile(`_test/event_attributes.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateEvent(test.event)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateAggregate(c *C) {

	aggregate := &cqrs.AggregateExpr{
		Name:   "item",
		Domain: s.domain,
		Commands: []*cqrs.CommandExpr{
			&cqrs.CommandExpr{
				Name:   "CreateItem",
				Events: []string{"ItemCreated"},
				Params: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "id",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.UUID,
								Validation: &gogen.ValidationRule{
									Required: true,
								},
							},
						},
						&gogen.Field{
							Name: "name",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
								Validation: &gogen.ValidationRule{
									Required: true,
								},
							},
						},
						&gogen.Field{
							Name: "test",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
								Validation: &gogen.ValidationRule{
									Required: false,
								},
							},
						},
					},
				},
			},
			&cqrs.CommandExpr{
				Name:   "UpdateItem",
				Events: []string{"ItemTitleUpdated", "ItemPriceUpdated"},
			},
			&cqrs.CommandExpr{
				Name:   "DeleteItem",
				Events: []string{"ItemDeleted"},
			},
		},
	}

	tests := []struct {
		name      string
		aggregate *cqrs.AggregateExpr
		expected  string
	}{
		{
			name:      "aggregate",
			aggregate: aggregate,
			expected:  loadfile(`./_test/aggregate.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateAggregate(test.aggregate)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateProjection(c *C) {

	tests := []struct {
		name       string
		projection *cqrs.ProjectionExpr
		expected   string
	}{
		{
			name: "projection",
			projection: &cqrs.ProjectionExpr{
				Name:         "items",
				HandlesEvent: []string{"ItemCreated", "ItemTitleUpdated", "ItemPriceUpdated", "ItemDeleted"},
				Domain:       s.domain,
			},
			expected: loadfile(`./_test/projection.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateProjection(test.projection)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateEventsFactory(c *C) {

	tests := []struct {
		name     string
		domain   *cqrs.DomainExpr
		expected string
	}{
		{
			name: "event factory 4 events",
			domain: &cqrs.DomainExpr{
				Name: "maindomain",
				Aggregates: []*cqrs.AggregateExpr{
					&cqrs.AggregateExpr{
						Name: "testAggregate1",
						Commands: []*cqrs.CommandExpr{
							&cqrs.CommandExpr{
								Name:   "UpdateItem",
								Events: []string{"ItemTitleUpdated", "ItemPriceUpdated"},
							},
							&cqrs.CommandExpr{
								Name:   "DeleteItem",
								Events: []string{"ItemDeleted"},
							},
						},
					},
					&cqrs.AggregateExpr{
						Name: "testAggregate2",
						Commands: []*cqrs.CommandExpr{
							&cqrs.CommandExpr{
								Name:   "CreateTest2",
								Events: []string{"Test2Created"},
							},
						},
					},
				},
				Events: map[string]*cqrs.EventExpr{
					"ItemCreated": &cqrs.EventExpr{ //<-- should be omitted
						Name: "ItemCreated",
					},
					"ItemPriceUpdated": &cqrs.EventExpr{
						Name: "ItemPriceUpdated",
					},
					"ItemTitleUpdated": &cqrs.EventExpr{
						Name: "ItemTitleUpdated",
					},
					"ItemDeleted": &cqrs.EventExpr{
						Name: "ItemDeleted",
					},
					"Test2Created": &cqrs.EventExpr{
						Name: "Test2Created",
					},
				},
			},
			expected: loadfile(`./_test/event_factory.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateEventsFactory(test.domain)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateAggregatesFactory(c *C) {

	tests := []struct {
		name     string
		domain   *cqrs.DomainExpr
		expected string
	}{
		{
			name: "2 aggregate factories",
			domain: &cqrs.DomainExpr{
				Name: "maindomain",
				Aggregates: []*cqrs.AggregateExpr{
					&cqrs.AggregateExpr{
						Name: "testAggregate1",
					},
					&cqrs.AggregateExpr{
						Name: "testAggregate2",
					},
				},
			},
			expected: loadfile(`./_test/aggregate_factory.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateAggregatesFactory(test.domain)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateRepositoryInterfaces(c *C) {

	tests := []struct {
		name     string
		domain   *cqrs.DomainExpr
		expected string
	}{
		{
			name: "1 repo no filter and 1 repo with filter definition",
			domain: &cqrs.DomainExpr{
				Projections: []*cqrs.ProjectionExpr{
					{
						ReadRepositories: []*cqrs.RepositoryExpr{
							{
								Name: "item",
								Model: &gogen.UserTypeExpr{
									TypeName: "testmodel1",
								},
							},
						},
					},
					{
						ReadRepositories: []*cqrs.RepositoryExpr{
							{
								Name: "product",
								Filter: &gogen.AttributeExpr{
									Type: &gogen.Object{
										&gogen.Field{
											Name: "id",
											Attribute: &gogen.AttributeExpr{
												Type: &gogen.Array{
													ElemType: &gogen.AttributeExpr{
														Type: gogen.UUID,
													},
												},
											},
										},
										&gogen.Field{
											Name: "name",
											Attribute: &gogen.AttributeExpr{
												Type: gogen.String,
											},
										},
										&gogen.Field{
											Name: "test",
											Attribute: &gogen.AttributeExpr{
												Type: &gogen.Array{
													ElemType: &gogen.AttributeExpr{
														Type: gogen.String,
													},
												},
											},
										},
									},
								},
								Model: &gogen.UserTypeExpr{
									TypeName: "testmodel2",
								},
							},
						},
					},
				},
			},
			expected: loadfile(`./_test/readrepository.go`),
		},
	}

	for _, test := range tests {
		sec, err := s.generator.GenerateRepositoryInterfaces(test.domain)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func (s *GeneratorSuite) TestGenerateDatabaseRepository(c *C) {

	model := &gogen.UserTypeExpr{
		TypeName: "item",
		AttributeExpr: &gogen.AttributeExpr{
			Type: &gogen.Object{
				&gogen.Field{
					Name: "Id",
					Attribute: &gogen.AttributeExpr{
						Type: gogen.UUID,
					},
				},
				&gogen.Field{
					Name: "Name",
					Attribute: &gogen.AttributeExpr{
						Type: gogen.String,
					},
				},
				&gogen.Field{
					Name: "Price",
					Attribute: &gogen.AttributeExpr{
						Type: gogen.Float32,
					},
				},
			},
		},
	}

	tests := []struct {
		name       string
		repository *cqrs.RepositoryExpr
		expected   string
	}{
		{
			name: "repo without filter",
			repository: &cqrs.RepositoryExpr{
				Name:  "productitem",
				Model: model,
			},
			expected: loadfile(`./_test/dbreadrepository_no_filter.go`),
		},
		{
			name: "repo with filter",
			repository: &cqrs.RepositoryExpr{
				Name: "productitem",
				Filter: &gogen.AttributeExpr{
					Type: &gogen.Object{
						&gogen.Field{
							Name: "id",
							Attribute: &gogen.AttributeExpr{
								Type: &gogen.Array{
									ElemType: &gogen.AttributeExpr{
										Type: gogen.UUID,
									},
								},
							},
						},
						&gogen.Field{
							Name: "name",
							Attribute: &gogen.AttributeExpr{
								Type: gogen.String,
							},
						},
						&gogen.Field{
							Name: "test",
							Attribute: &gogen.AttributeExpr{
								Type: &gogen.Array{
									ElemType: &gogen.AttributeExpr{
										Type: gogen.String,
									},
								},
							},
						},
					},
				},
				Model: model,
			},
			expected: loadfile(`./_test/dbreadrepository.go`),
		},


	}

	for _, test := range tests {
		sec, err := s.generator.GenerateDbRepository(test.repository)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func loadfile(filename string) string {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(buf, f)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
