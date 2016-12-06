package memory

import (
	"github.com/looplab/eventhorizon"
	. "gopkg.in/check.v1"
	"github.com/mbict/go-cqrs"
	"testing"
)

type testEventA struct {
	*cqrs.EventBase
	A int
	B string
}

func (e *testEventA) EventType() string { return "testEventA" }

type testEventB struct {
	*cqrs.EventBase
	C string
	D bool
}

func (e *testEventB) EventType() string { return "testEventB" }

func Test(t *testing.T) { TestingT(t) }
