package cqrs

import (
	"fmt"
	"strconv"
)

type AggregateId interface {
	// String gives the textual representation back of the AggregateId
	String() string

	//Value should give back the real used type back
	Value() interface{}
}

type stringerAggregateId struct {
	fmt.Stringer
}

func (s *stringerAggregateId) Value() interface{} {
	return s.Stringer
}

func NewAggregateId(id fmt.Stringer) AggregateId {
	return &stringerAggregateId{Stringer: id}
}

type stringAggregateId struct {
	id string
}

func (s *stringAggregateId) String() string {
	return s.id
}

func (s *stringAggregateId) Value() interface{} {
	return s.id
}

func NewStringAggregateId(id string) AggregateId {
	return &stringAggregateId{id: id}
}

type intAggregateId struct {
	id int
}

func (i *intAggregateId) String() string {
	return strconv.Itoa(i.id)
}

func (i *intAggregateId) Value() interface{} {
	return i.id
}

func NewIntAggregateId(id int) AggregateId {
	return &intAggregateId{id: id}
}
