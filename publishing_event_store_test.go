package cqrs

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestShouldWriteEventsAndPublishEvents(t *testing.T) {
	es := &MockEventStore{}
	es.On("WriteEvent", mock.Anything, mock.Anything).Return(nil)
	var publishCalledWith []Event
	fn := EventPublisherFunc(func(events ...Event) error {
		publishCalledWith = append(publishCalledWith, events...)
		return nil
	})
	eventData := &eventA{}
	event := NewEvent(uuid.Must(uuid.NewV4()), 1, time.Now(), eventData)
	pes := NewEventPublishingEventStore(fn, es)

	err := pes.WriteEvent("test", event)

	assert.NoError(t, err, "publish should not return an error")
	assert.Equal(t, 1, len(publishCalledWith), "times publish callback will be called")

	es.AssertNumberOfCalls(t, "WriteEvent", 1)
	es.AssertCalled(t, "WriteEvent", "test", event)
}

func TestShouldNotPublishEventsOnEventstoreWriteError(t *testing.T) {
	expectedError := errors.New("error")
	es := &MockEventStore{}
	es.On("WriteEvent", mock.Anything, mock.Anything).Return(expectedError)
	var publishCalledWith []Event
	fn := EventPublisherFunc(func(events ...Event) error {
		publishCalledWith = append(publishCalledWith, events...)
		return nil
	})
	eventData := &eventA{}
	event := NewEvent(uuid.Must(uuid.NewV4()), 1, time.Now(), eventData)
	pes := NewEventPublishingEventStore(fn, es)

	err := pes.WriteEvent("test", event)

	assert.Error(t, err, expectedError)
	assert.Equal(t, 0, len(publishCalledWith), "publish callback should not be called")

	es.AssertNumberOfCalls(t, "WriteEvent", 1)
	es.AssertCalled(t, "WriteEvent", "test", event)
}

func TestShouldWriteEventsAndPublishEventsAndReturnErrorFromPublish(t *testing.T) {
	expectedError := errors.New("error")
	es := &MockEventStore{}
	es.On("WriteEvent", mock.Anything, mock.Anything).Return(nil)
	var publishCalledWith []Event
	fn := EventPublisherFunc(func(events ...Event) error {
		publishCalledWith = append(publishCalledWith, events...)
		return expectedError
	})
	eventData := &eventA{}
	event := NewEvent(uuid.Must(uuid.NewV4()), 1, time.Now(), eventData)
	pes := NewEventPublishingEventStore(fn, es)

	err := pes.WriteEvent("test", event)

	assert.Error(t, err, expectedError)
	assert.Equal(t, 1, len(publishCalledWith), "times publish callback will be called")

	es.AssertNumberOfCalls(t, "WriteEvent", 1)
	es.AssertCalled(t, "WriteEvent", "test", event)
}
