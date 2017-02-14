package cqrs

import (
	"github.com/mbict/gogen"
)

type RepositoryExpr struct {
	//Name of the read repository
	Name string

	// Model this repository wil provide
	Model *gogen.UserTypeExpr

	// Filter defines the fields of the filter for the find function in the repository
	Filter *gogen.AttributeExpr
}

func (r *RepositoryExpr) Attribute() *gogen.AttributeExpr {
	return r.Filter
}

func (r *RepositoryExpr) Context() string {
	return "repository"
}
