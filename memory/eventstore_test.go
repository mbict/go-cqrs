package memory

import (
	. "gopkg.in/check.v1"
	"github.com/mbict/go-cqrs"
)

var _ = Suite(&EventStoreSuite{})

type EventStoreSuite struct {
	store cqrs.EventStore
}

func (s *EventStoreSuite) SetUpTest(c *C) {
	s.store = NewMemoryEventStore()
}

func (s *EventStoreSuite) TestLoadStream(c *C) {

	//c.Assert(eh.events, HasLen, 1)
	//c.Assert(eh.events[0], DeepEquals, e)
}

func (s *EventStoreSuite) TestWrite(c *C) {
	s.store.WriteEvent("testAggregate", cqrs.NewEventMessage("100", &testEventA{A: 101, B: "foo"}, 1))
	s.store.WriteEvent("testAggregate", cqrs.NewEventMessage("100", &testEventB{C: "bar", D: true}, 2))
	s.store.WriteEvent("testAggregate", cqrs.NewEventMessage("101", &testEventA{A: 102, B: "foo"}, 1))
	s.store.WriteEvent("testAggregate", cqrs.NewEventMessage("100", &testEventB{C: "foo", D: false}, 3))
	s.store.WriteEvent("testAggregate2", cqrs.NewEventMessage("100", &testEventA{A: 103, B: "foo"}, 1))

	stream, err := s.store.LoadStream("testAggregate", "100")

	eventA := testEventA{}
	eventB := testEventB{}

	c.Assert(err, IsNil)

	c.Assert(stream.Next(), Equals, true)
	c.Assert(stream.EventType(), Equals, "testEventA")
	c.Assert(stream.Version(), Equals, 1)
	c.Assert(stream.Scan(&eventA), IsNil)
	c.Assert(eventA.A, Equals, 101)
	c.Assert(eventA.B, Equals, "foo")

	c.Assert(stream.Next(), Equals, true)
	c.Assert(stream.EventType(), Equals, "testEventB")
	c.Assert(stream.Version(), Equals, 2)
	c.Assert(stream.Scan(&eventB), IsNil)
	c.Assert(eventB.C, Equals, "bar")
	c.Assert(eventB.D, Equals, true)

	c.Assert(stream.Next(), Equals, true)
	c.Assert(stream.EventType(), Equals, "testEventB")
	c.Assert(stream.Version(), Equals, 3)
	c.Assert(stream.Scan(&eventB), IsNil)
	c.Assert(eventB.C, Equals, "foo")
	c.Assert(eventB.D, Equals, false)

	c.Assert(stream.Next(), Equals, false)
	c.Assert(stream.EventType(), Equals, "")
	c.Assert(stream.Version(), Equals, -1)
	c.Assert(stream.Scan(&eventB), Equals, ErrNoEventData)
}
