package memory

import (
	. "gopkg.in/check.v1"
	"github.com/mbict/go-cqrs"
)

type CqrsSuite struct {
	dispatcher       *cqrs.Dispatcher
	aggregateFactory *cqrs.CallbackAggregateFactory
	eventFactory     *cqrs.CallbackEventFactory
	eventStore       cqrs.EventStore
	repository       cqrs.DomainRepository
}

func (s *CqrsSuite) SetUpTest(c *C) {
	s.eventStore = NewMemoryEventStore()
	s.eventFactory = cqrs.NewCallbackEventFactory()
	s.aggregateFactory = cqrs.NewCallbackAggregateFactory()
	s.repository = cqrs.NewCommonDomainRepository(s.eventStore, s.eventFactory, s.aggregateFactory)
	s.dispatcher, _ = cqrs.NewDispatcher(s.repository)
}

func (s *CqrsSuite) TestPublishEventsToHandlers(c *C) {
	s.aggregateFactory.RegisterCallback(func(id string) cqrs.AggregateRoot {
		return &testAggregate{AggregateBase: cqrs.NewAggregateBase(id)}
	})

	s.eventFactory.RegisterCallback(func(id string) cqrs.Event {
		return &testEventA{
			EventBase: cqrs.NewEventBase(id),
		}
	})

	s.eventFactory.RegisterCallback(func(id string) cqrs.Event {
		return &testEventB{
			EventBase: cqrs.NewEventBase(id),
		}
	})
}

type testAggregate struct {
	*cqrs.AggregateBase
}

func (a *testAggregate) Apply(cqrs.Event) error {
	return nil
}

func (a *testAggregate) AggregateType() string {
	return "testAggregate"
}

func (a *testAggregate) HandleCommand(cqrs.Command) error {
	return nil
}
