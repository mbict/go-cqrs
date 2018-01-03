package cqrs

import (
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewEventBaseFromAggregate(t *testing.T) {
	id := uuid.NewV4()
	agg := &MockAggregate{}
	agg.On("AggregateId").Return(id)
	agg.On("CurrentVersion").Return(10)

	eventBase := NewEventBaseFromAggregate(agg)

	if eventBase.Version() != 11 {
		t.Errorf("expected version %d but got %d", 11, eventBase.Version())
	}

	if !uuid.Equal(eventBase.AggregateId(), id) {
		t.Errorf("expected aggregate id `%s` but got `%s`", id.String(), eventBase.AggregateId())
	}

	if eventBase.OccurredAt().IsZero() {
		t.Error("expected aggregate occurred at to be arround now but got an empty time")
	}
}
