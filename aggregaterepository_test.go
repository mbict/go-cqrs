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
	occurredAt := time.Now()
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("OccurredAt").Return(occurredAt)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(1)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1, occurredAt).Return(nil)

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
	occurredAt := time.Now()
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(999)
	stream.On("OccurredAt").Return(occurredAt)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0)
	aggregateFactory := func(ctx AggregateContext) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 999, occurredAt).Return(event)

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
	occurredAt := time.Now()
	scanError := errors.New("scan error")
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("OccurredAt").Return(occurredAt)
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
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1, occurredAt).Return(event)

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
	stream.On("EventName").Return("")
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
	occurredAt := time.Now()
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Once().Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("OccurredAt").Return(occurredAt)
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
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1, occurredAt).Return(event)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg == nil {
		t.Errorf("expected non nil aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 2)
	stream.AssertNumberOfCalls(t, "Scan", 1)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertCalled(t, "Apply", *event)
	//	aggregate.AssertNumberOfCalls(t, "incrementVersion", 1)
}

func TestAggregateRepository_LoadWithMultipleEvents(t *testing.T) {
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Times(3)
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Times(3).Return("testEvent")
	stream.On("OccurredAt").Return(time.Now()).Times(3)
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
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, mock.Anything, mock.Anything).Return(event)

	repo := NewAggregateRepository(store, DefaultAggregateBuilder(aggregateFactory), eventFactory)

	agg, err := repo.Load(uuid.Must(uuid.NewV4()))
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
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 3)
	aggregate.AssertCalled(t, "Apply", *event)
}

/*********************
  Save tests
 *********************/

func TestAggregateRepository_SaveWithErrorWriteEvent(t *testing.T) {
	storeError := errors.New("store write error")
	event := &eventA{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)

	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event).Return(storeError)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != storeError {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
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
	event := &eventA{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)
	aggregate.On("clearUncommittedEvents")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 1)
}

func TestAggregateRepository_SaveWithMultipleEvents(t *testing.T) {
	event := &eventA{}
	events := []Event{event, *event, event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)
	aggregate.On("clearUncommittedEvents")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event, *event, event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event, *event, event)
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 1)
}

func TestAggregateRepository_SaveTriggersPublishEventHook(t *testing.T) {
	hookCalled := 0
	var hookArgument Event
	eventPublishHook := func(event Event) {
		hookCalled++
		hookArgument = event
	}

	event := &eventA{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("getUncommittedEvents").Return(events)
	aggregate.On("clearUncommittedEvents")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil, eventPublishHook)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if hookCalled != 1 {
		t.Errorf("expected event publish hook to be called %d times, but got called %d times", 1, hookCalled)
	}

	if hookArgument != *event {
		t.Errorf("expected event publish hook to have event `%v` as argument but got `%v`", event, hookArgument)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
	aggregate.AssertNumberOfCalls(t, "getUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "clearUncommittedEvents", 1)
}
