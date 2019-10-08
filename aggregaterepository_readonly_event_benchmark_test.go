package cqrs

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

//simple benchmark tests to see if the reflect will impact real world usage

var fixedId = uuid.New()

type test_eventStream struct {
	i         int
	count     int
	timestamp time.Time
}

func (s *test_eventStream) EventType() EventType {
	return ""
}

func (s *test_eventStream) Scan(EventData) error {
	return nil
}

func (*test_eventStream) EventName() string {
	return "event.a"
}

func (*test_eventStream) AggregateId() AggregateId {
	return fixedId
}

func (s *test_eventStream) Version() int {
	return s.i
}

func (s *test_eventStream) Timestamp() time.Time {
	return s.timestamp
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

var x = 1

func callMe(e EventData) {
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
	id := uuid.New()
	event := &eventA{}

	eventStore := &MockEventStore{}
	eventStore.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(&test_eventStream{i: 0, count: 100}, nil)
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", mock.Anything, mock.Anything, mock.Anything).Return(event, nil)
	repository := NewAggregateRepository(eventStore, DefaultAggregateBuilder(aggregateAFactory), eventFactory)

	for n := 0; n < b.N; n++ {
		//eventStream := &MockEventStream{}
		//eventStream.On("Next").Return(true).Times(1000)
		//eventStream.On("Next").Return(false)
		//eventStream.On("EventType").Return("event.a")
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
//	id := uuid.Must(uuid.NewV4())
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
//		//eventStream.On("EventType").Return("event.a")
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
