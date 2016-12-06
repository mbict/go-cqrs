package cqrs

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&LocalEventBusSuite{})

type LocalEventBusSuite struct {
	bus EventBus
}

func (s *LocalEventBusSuite) SetUpTest(c *C) {
	s.bus = NewLocalEventBus()
}

func (s *LocalEventBusSuite) TestPublishEventsToHandlers(c *C) {
	eh := NewMockEventHandler()
	s.bus.AddHandler(eh, &testEvent{})
	e := NewEventMessage("1", &testEvent{}, 1)

	s.bus.PublishEvent(e)
	s.bus.PublishEvent(NewEventMessage("2", &testOtherEvent{}, 1))

	c.Assert(eh.events, HasLen, 1)
	c.Assert(eh.events[0], DeepEquals, e)
}

func (s *LocalEventBusSuite) TestAddHandlerAndIgnoreDuplicateHandlers(c *C) {
	e1 := NewEventMessage("1", &testEvent{}, 1)
	e2 := NewEventMessage("1", &testOtherEvent{}, 1)
	eh := NewMockEventHandler()
	s.bus.AddHandler(eh, &testEvent{})
	s.bus.AddHandler(eh, &testEvent{}, &testOtherEvent{})

	s.bus.PublishEvent(e1)
	s.bus.PublishEvent(e2)

	c.Assert(eh.events, HasLen, 2)
	c.Assert(eh.events[0], DeepEquals, e1)
	c.Assert(eh.events[1], DeepEquals, e2)

}
