package cqrs

import (
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewEventBaseFromAggregate(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	agg := &MockAggregateContext{}
	agg.On("AggregateId").Return(id)
	agg.On("OriginalVersion").Return(10)
	agg.On("getUncommittedEvents").Return(nil)

	event := NewEventFromAggregate(agg, eventA{})

	if event.Version() != 11 {
		t.Errorf("expected version %d but got %d", 11, event.Version())
	}

	if string(event.EventType()) != "event:a" {
		t.Errorf("expected event type to be `%s` but got `%s`", "event:a", event.EventType())
	}

	if !uuid.Equal(event.AggregateId(), id) {
		t.Errorf("expected aggregate id `%s` but got `%s`", id.String(), event.AggregateId())
	}

	if event.Timestamp().IsZero() {
		t.Error("expected aggregate occurred at to be arround now but got an empty time")
	}
}
