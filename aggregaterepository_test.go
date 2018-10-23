package cqrs

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

/*********************
  Load tests
 *********************/

func TestAggregateRepository_LoadWithStreamError(t *testing.T) {
	streamError := errors.New("stream error")
	stream := &MockEventStream{}
	stream.On("Next").Return(false)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(nil, streamError)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(1)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), nil)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err == nil || err.Error() != "cannot load events from stream reader, error: stream error" {
		t.Errorf("expected error `%v`, but got error `%v`", streamError, err)
	}

	if agg != nil {
		t.Errorf("expected a nil aggregate , but got `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 0)
}

func TestAggregateRepository_LoadWithUnkownEventFactoryError(t *testing.T) {
	timestamp := time.Now()
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType("testEvent"))
	stream.On("Version").Return(2)
	stream.On("Timestamp").Return(timestamp)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(1)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", EventType("testEvent")).Return(nil)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err == nil || err.Error() != "the repository has no event factory registered for event type: testEvent" {
		t.Errorf("expected error `%v`, but got error `%v`", "", err)
	}

	if agg != nil {
		t.Errorf("expected a nil aggregate , but got `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
}

func TestAggregateRepository_LoadWithVersionMismatch(t *testing.T) {
	timestamp := time.Now()
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType("testEvent"))
	stream.On("Version").Return(999)
	stream.On("Timestamp").Return(timestamp)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", EventType("testEvent")).Return(event)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err == nil || err.Error() != "event version (999) mismatch with Aggregate next Version (1)" {
		t.Errorf("expected version mismatch error `%v`, but got error `%v`", "event version (999) mismatch with Aggregate next Version (1)", err)
	}

	if agg != nil {
		t.Errorf("expected a nil aggregate , but got `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
}

func TestAggregateRepository_LoadWithScanFailure(t *testing.T) {
	timestamp := time.Now()
	scanError := errors.New("scan error")
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType("testEvent"))
	stream.On("Version").Return(1)
	stream.On("Timestamp").Return(timestamp)
	stream.On("Scan", mock.Anything).Return(scanError)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("Version").Return(0)
	aggregate.On("AggregateName").Return("testAggregate")
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", EventType("testEvent")).Return(event)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err == nil || err.Error() != "the repository cannot populate event data from stream for event type: testEvent, with error `scan error`" {
		t.Errorf("expected a scan error, but got error `%v`", err)
	}

	if agg != nil {
		t.Errorf("expected a nil aggregate , but got `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
	stream.AssertNumberOfCalls(t, "Scan", 1)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 0)
	aggregate.AssertNumberOfCalls(t, "incrementVersion", 0)
}

// success paths
func TestAggregateRepository_LoadWithNoEvents(t *testing.T) {
	stream := &MockEventStream{}
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType(""))
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(1)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), nil)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	_, ok := agg.(Aggregate)
	if agg == nil || !ok {
		t.Errorf("expected a non nil aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
	stream.AssertNumberOfCalls(t, "Scan", 0)
	aggregate.AssertNumberOfCalls(t, "Apply", 0)
	aggregate.AssertNumberOfCalls(t, "incrementVersion", 0)
}

func TestAggregateRepository_LoadWithOneEvent(t *testing.T) {
	aggregateId := uuid.Must(uuid.NewV4())
	timestamp := time.Now()
	eventData := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType("")).Once().Return(EventType("testEvent"))
	stream.On("Version").Return(1)
	stream.On("Timestamp").Return(timestamp)
	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0)
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("incrementVersion")
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", EventType("testEvent")).Return(eventData)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(aggregateId)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg == nil {
		t.Errorf("expected non nil aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 2)
	stream.AssertNumberOfCalls(t, "Scan", 1)
	stream.AssertCalled(t, "Scan", eventData)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertNumberOfCalls(t, "incrementVersion", 1)
	aggregate.AssertCalled(t, "Apply", matchEventData(aggregateId, *eventData, 1))
}

func TestAggregateRepository_LoadWithMultipleEvents(t *testing.T) {
	eventData := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Times(3)
	stream.On("Next").Return(false)
	stream.On("EventType").Return(EventType("")).Times(3).Return(EventType("testEvent"))
	stream.On("Timestamp").Return(time.Now()).Times(3)
	stream.On("Version").Return(1).Times(2)
	stream.On("Version").Return(2).Times(2)
	stream.On("Version").Return(3).Times(2)
	stream.On("Version").Return(0)

	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0).Once() //load
	aggregate.On("Version").Return(0).Once() //check version
	aggregate.On("Version").Return(1).Once() //check version
	aggregate.On("Version").Return(2).Once() //check version
	aggregate.On("Version").Return(3)        //final version
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("incrementVersion")
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", EventType("testEvent")).Return(eventData)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	aggregateId := uuid.Must(uuid.NewV4())
	agg, err := repo.Load(aggregateId)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	_, ok := agg.(Aggregate)
	if agg == nil || !ok {
		t.Errorf("expected non nil aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 4)
	stream.AssertNumberOfCalls(t, "Scan", 3)
	stream.AssertCalled(t, "Scan", eventData)
	aggregate.AssertNumberOfCalls(t, "Apply", 3)
	aggregate.AssertCalled(t, "Apply", matchEventData(aggregateId, *eventData, 1))
}

/*********************
  Save tests
 *********************/

func TestAggregateRepository_SaveWithErrorWriteEvent(t *testing.T) {
	aggregateId := uuid.Must(uuid.NewV4())
	storeError := errors.New("store write error")
	eventData := &eventA{}
	events := []Event{NewEvent(aggregateId, 1, time.Now(), eventData)}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)

	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", mock.Anything).Return(storeError)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != storeError {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", matchEventData(aggregateId, *eventData, 1))
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "Apply", 0)
	aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 0)
}

// success paths
func TestAggregateRepository_SaveWithNoEvents(t *testing.T) {
	aggregate := &MockAggregate{}
	aggregate.On("getUncommittedEvents").Return(nil)
	store := &MockEventStore{}

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 0)
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
}

func TestAggregateRepository_SaveWithOneEvent(t *testing.T) {
	eventData := &eventA{}
	aggregateId := uuid.Must(uuid.NewV4())
	events := []Event{NewEvent(aggregateId, 1, time.Now(), eventData)}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)
	aggregate.On("clearUncommittedEvents")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", mock.Anything).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	//store.AssertNumberOfCalls(t, "WriteEvent", 1)
	//store.AssertCalled(t, "WriteEvent", "testAggregate", matchEventData( aggregateId, eventData, 1))
	//aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	//aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 1)
}

func TestAggregateRepository_SaveWithMultipleEvents(t *testing.T) {
	event := NewEvent(uuid.Must(uuid.NewV4()), 1, time.Now(), &eventA{})
	events := []Event{event, event, event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)
	aggregate.On("clearUncommittedEvents")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event, event, event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event, event, event)
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 1)
}

func matchEventData(aggregateId uuid.UUID, eventData EventData, version int) interface{} {
	return mock.MatchedBy(func(e Event) bool {
		return e != nil &&
			e.Data() != nil &&
			uuid.Equal(e.AggregateId(), aggregateId) &&
			e.EventType() == "event:a" &&
			e.Data() == eventData &&
			e.Version() == version
	})
}
