package cqrs

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

//simple benchmark tests to see if the reflect will impact real world usage

var fixedId = uuid.NewV4()

type test_eventStore struct{}

func (*test_eventStore) LoadStream(aggregateName string, aggregateId uuid.UUID) (EventStream, error) {
	return &test_eventStream{i: 0, count: 100}, nil
}

func (*test_eventStore) WriteEvent(string, ...Event) error {
	panic("implement me")
}

type test_eventStream struct {
	i     int
	count int
}

func (*test_eventStream) EventName() string {
	return "event.a"
}

func (*test_eventStream) AggregateId() uuid.UUID {
	return fixedId
}

func (s *test_eventStream) Version() int {
	return s.i
}

func (s *test_eventStream) Next() bool {
	if s.i <= s.count {
		s.i++
		return true
	}
	s.i = 0
	return false
}

func (*test_eventStream) Error() error {
	return nil
}

func (*test_eventStream) Scan(Event) error {
	return nil
}

var x = 1

func callMe(e Event) {
	x++
}

func BenchmarkDirectPass(b *testing.B) {
	for n := 0; n < b.N; n++ {
		event := &eventA{}
		callMe(event)
	}
}

func BenchmarkReflectOverhead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		event := &eventA{}
		callMe(reflect.Indirect(reflect.ValueOf(event)).Interface().(Event))
	}
}

func BenchmarkAggregateRepositoryEvenPassedByValue(b *testing.B) {
	id := uuid.NewV4()
	event := &eventA{}

	eventStore := &test_eventStore{}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", mock.Anything, mock.Anything, mock.Anything).Return(event, nil)
	repository := NewAggregateRepository(eventStore, aggregateAFactory, eventFactory)

	for n := 0; n < b.N; n++ {
		//eventStream := &MockEventStream{}
		//eventStream.On("Next").Return(true).Times(1000)
		//eventStream.On("Next").Return(false)
		//eventStream.On("EventName").Return("event.a")
		//eventStream.On("Version").Return(1)
		//eventStream.On("Scan", mock.Anything).Return(nil)
		//
		//eventStore := &MockEventStore{}
		//eventStore.On("LoadStream", mock.Anything, mock.Anything).Return(eventStream, nil)
		//eventFactory := &MockEventFactory{}
		//eventFactory.On("MakeEvent", mock.Anything, mock.Anything, mock.Anything).Return(event, nil)
		//repository := NewAggregateRepositoryReadonlyEvents(eventStore, aggregateAFactory, eventFactory)

		repository.Load(id)
	}
}

//func BenchmarkAggregateRepositoryEvenPassedByReference(b *testing.B) {
//	id := uuid.NewV4()
//	event := &eventA{}
//
//	eventStore := &test_eventStore{}
//	eventFactory := &MockEventFactory{}
//	eventFactory.On("MakeEvent", mock.Anything, mock.Anything, mock.Anything).Return(event, nil)
//	repository := NewAggregateRepositoryReadonlyEvents(eventStore, aggregateAFactory, eventFactory)
//
//	for n := 0; n < b.N; n++ {
//		//eventStream := &MockEventStream{}
//		//eventStream.On("Next").Return(true).Times(1000)
//		//eventStream.On("Next").Return(false)
//		//eventStream.On("EventName").Return("event.a")
//		//eventStream.On("Version").Return(1)
//		//eventStream.On("Scan", mock.Anything).Return(nil)
//		//
//		//eventStore := &MockEventStore{}
//		//eventStore.On("LoadStream", mock.Anything, mock.Anything).Return(eventStream, nil)
//		//eventFactory := &MockEventFactory{}
//		//eventFactory.On("MakeEvent", mock.Anything, mock.Anything, mock.Anything).Return(event, nil)
//		//repository := NewAggregateRepository(eventStore, aggregateAFactory, eventFactory)
//
//		repository.Load(id)
//	}
//}
