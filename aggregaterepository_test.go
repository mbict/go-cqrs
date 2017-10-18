package cqrs

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*********************
  Load tests
 *********************/

func TestAggregateRepository_LoadWithStreamError(t *testing.T) {
	streamError := errors.New("stream error")
	stream := &MockEventStream{}
	stream.On("Next").Return(false)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(nil, streamError)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	repo := NewAggregateRepository(store, aggregateFactory, nil)

	agg, err := repo.Load(uuid.NewV4())
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
	//event := &MockEvent{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(1)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1).Return(nil)

	repo := NewAggregateRepository(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
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
	event := &MockEvent{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(9999999)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(11111111)
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 9999999).Return(event)

	repo := NewAggregateRepository(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
	if err == nil || err.Error() != "event version (9999999) mismatch with Aggregate next Version (11111112)" {
		t.Errorf("expected version mismatch error `%v`, but got error `%v`", "event version (9999999) mismatch with Aggregate next Version (11111112)", err)
	}

	if agg != nil {
		t.Errorf("expected a nil aggregate , but got `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
}

func TestAggregateRepository_LoadWithScanFailure(t *testing.T) {
	scanError := errors.New("scan error")
	event := &MockEvent{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("Scan", mock.Anything).Return(scanError)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("Version").Return(0)
	aggregate.On("AggregateName").Return("testAggregate")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1).Return(event)

	repo := NewAggregateRepository(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
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
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 0)
}

// success paths
func TestAggregateRepository_LoadWithNoEvents(t *testing.T) {
	stream := &MockEventStream{}
	stream.On("Next").Return(false)
	stream.On("EventName").Return("")
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}

	repo := NewAggregateRepository(store, aggregateFactory, nil)

	agg, err := repo.Load(uuid.NewV4())
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg != aggregate {
		t.Errorf("expected identical aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 1)
	stream.AssertNumberOfCalls(t, "Scan", 0)
	aggregate.AssertNumberOfCalls(t, "Apply", 0)
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 0)
}

func TestAggregateRepository_LoadWithOneEvent(t *testing.T) {
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Once().Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0)
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("IncrementVersion")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1).Return(event)

	repo := NewAggregateRepository(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg != aggregate {
		t.Errorf("expected identical aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 2)
	stream.AssertNumberOfCalls(t, "Scan", 1)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertCalled(t, "Apply", event)
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 1)
}

func TestAggregateRepository_LoadWithMultipleEvents(t *testing.T) {
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Times(3)
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Times(3).Return("testEvent")
	stream.On("Version").Return(1).Times(2)
	stream.On("Version").Return(2).Times(2)
	stream.On("Version").Return(3).Times(2)
	stream.On("Version").Return(0)

	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0).Once()
	aggregate.On("Version").Return(1).Once()
	aggregate.On("Version").Return(2).Once()
	aggregate.On("Version").Return(3)
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("IncrementVersion")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, mock.Anything).Return(event)

	repo := NewAggregateRepository(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg != aggregate {
		t.Errorf("expected identical aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 4)
	stream.AssertNumberOfCalls(t, "Scan", 3)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 3)
	aggregate.AssertCalled(t, "Apply", event)
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 3)
}

func TestAggregateRepository_LoadWithOneEventReadonlyEvents(t *testing.T) {
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Once()
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Once().Return("testEvent")
	stream.On("Version").Return(1)
	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0)
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("IncrementVersion")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, 1).Return(event)

	repo := NewAggregateRepositoryReadonlyEvents(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg != aggregate {
		t.Errorf("expected identical aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 2)
	stream.AssertNumberOfCalls(t, "Scan", 1)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertCalled(t, "Apply", *event)
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 1)
}

func TestAggregateRepository_LoadWithMultipleEventsReadonlyEvents(t *testing.T) {
	event := &eventA{}
	stream := &MockEventStream{}
	stream.On("Next").Return(true).Times(3)
	stream.On("Next").Return(false)
	stream.On("EventName").Return("").Times(3).Return("testEvent")
	stream.On("Version").Return(1).Times(2)
	stream.On("Version").Return(2).Times(2)
	stream.On("Version").Return(3).Times(2)
	stream.On("Version").Return(0)

	stream.On("Scan", mock.Anything).Return(nil)
	store := &MockEventStore{}
	store.On("LoadStream", mock.Anything, mock.Anything).Return(stream, nil)
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("Version").Return(0).Once()
	aggregate.On("Version").Return(1).Once()
	aggregate.On("Version").Return(2).Once()
	aggregate.On("Version").Return(3)
	aggregate.On("Apply", mock.Anything).Return(nil)
	aggregate.On("IncrementVersion")
	aggregateFactory := func(id uuid.UUID) Aggregate {
		return aggregate
	}
	eventFactory := &MockEventFactory{}
	eventFactory.On("MakeEvent", "testEvent", mock.Anything, mock.Anything).Return(event)

	repo := NewAggregateRepositoryReadonlyEvents(store, aggregateFactory, eventFactory)

	agg, err := repo.Load(uuid.NewV4())
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	if agg != aggregate {
		t.Errorf("expected identical aggregate, but got aggregate `%v`", agg)
	}

	store.AssertNumberOfCalls(t, "LoadStream", 1)
	stream.AssertNumberOfCalls(t, "Next", 4)
	stream.AssertNumberOfCalls(t, "Scan", 3)
	stream.AssertCalled(t, "Scan", event)
	aggregate.AssertNumberOfCalls(t, "Apply", 3)
	aggregate.AssertCalled(t, "Apply", *event)
	aggregate.AssertNumberOfCalls(t, "IncrementVersion", 3)
}

/*********************
  Save tests
 *********************/

func TestAggregateRepository_SaveWithErrorWriteEvent(t *testing.T) {
	storeError := errors.New("store write error")
	event := &MockEvent{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("GetUncommittedEvents").Return(events)
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event).Return(storeError)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != storeError {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
	aggregate.AssertNumberOfCalls(t, "GetUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "Apply", 0)
	aggregate.AssertNumberOfCalls(t, "ClearUncommittedEvents", 0)
}

// success paths
func TestAggregateRepository_SaveWithNoEvents(t *testing.T) {
	aggregate := &MockAggregate{}
	aggregate.On("GetUncommittedEvents").Return(nil)
	store := &MockEventStore{}

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 0)
	aggregate.AssertNumberOfCalls(t, "GetUncommittedEvents", 1)
}

func TestAggregateRepository_SaveWithOneEvent(t *testing.T) {
	event := &MockEvent{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("GetUncommittedEvents").Return(events)
	aggregate.On("ClearUncommittedEvents")
	aggregate.On("Apply", event).Return(nil)
	aggregate.On("IncrementVersion")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
	aggregate.AssertNumberOfCalls(t, "GetUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertCalled(t, "Apply", event)
	aggregate.AssertNumberOfCalls(t, "ClearUncommittedEvents", 1)
}

func TestAggregateRepository_SaveWithMultipleEvents(t *testing.T) {
	event := &MockEvent{}
	events := []Event{event, event, event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("GetUncommittedEvents").Return(events)
	aggregate.On("ClearUncommittedEvents")
	aggregate.On("Apply", event).Return(nil)
	aggregate.On("IncrementVersion")
	store := &MockEventStore{}
	store.On("WriteEvent", "testAggregate", event, event, event).Return(nil)

	repo := NewAggregateRepository(store, nil, nil)

	err := repo.Save(aggregate)
	if err != nil {
		t.Errorf("expected a nil error, but got error `%v`", err)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event, event, event)
	aggregate.AssertNumberOfCalls(t, "GetUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "Apply", 3)
	aggregate.AssertCalled(t, "Apply", event)
	aggregate.AssertNumberOfCalls(t, "ClearUncommittedEvents", 1)
}

func TestAggregateRepository_SaveTriggersPublishEventHook(t *testing.T) {
	hookCalled := 0
	var hookArgument Event
	eventPublishHook := func(event Event) {
		hookCalled++
		hookArgument = event
	}

	event := &MockEvent{}
	events := []Event{event}
	aggregate := &MockAggregate{}
	aggregate.On("AggregateName").Return("testAggregate")
	aggregate.On("GetUncommittedEvents").Return(events)
	aggregate.On("ClearUncommittedEvents")
	aggregate.On("Apply", event).Return(nil)
	aggregate.On("IncrementVersion")
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

	if hookArgument != event {
		t.Errorf("expected event publish hook to have event `%v` as argument but got `%v`", event, hookArgument)
	}

	store.AssertNumberOfCalls(t, "WriteEvent", 1)
	store.AssertCalled(t, "WriteEvent", "testAggregate", event)
	aggregate.AssertNumberOfCalls(t, "GetUncommittedEvents", 1)
	aggregate.AssertNumberOfCalls(t, "Apply", 1)
	aggregate.AssertCalled(t, "Apply", event)
	aggregate.AssertNumberOfCalls(t, "ClearUncommittedEvents", 1)
}
