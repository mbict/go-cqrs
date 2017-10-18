package cqrs

import (
	"github.com/satori/go.uuid"
	"testing"
)

func TestAggregateBase_AggregateId(t *testing.T) {
	id := uuid.NewV4()
	base := NewAggregateBase(id)

	if uuid.Equal(base.AggregateId(), id) == false {
		t.Errorf("expected aggregateId %s but got %s", id.String(), base.AggregateId().String())
	}
}

func TestAggregateBase_EventVersioning(t *testing.T) {
	id := uuid.NewV4()
	base := NewAggregateBase(id)

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.CurrentVersion())
	}

	base.StoreEvent(&MockEvent{})

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 1 {
		t.Errorf("expected version %d but got %d", 1, base.CurrentVersion())
	}
}

func TestAggregateBase_IncrementVersion(t *testing.T) {
	id := uuid.NewV4()
	base := NewAggregateBase(id)

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.CurrentVersion())
	}

	base.IncrementVersion()

	if base.Version() != 1 {
		t.Errorf("expected version %d but got %d", 1, base.Version())
	}

	if base.CurrentVersion() != 1 {
		t.Errorf("expected version %d but got %d", 1, base.CurrentVersion())
	}
}

func TestAggregateBase_EventsHandling(t *testing.T) {
	id := uuid.NewV4()
	base := NewAggregateBase(id)

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.CurrentVersion())
	}

	base.StoreEvent(&MockEvent{})
	base.StoreEvent(&MockEvent{})

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 2 {
		t.Errorf("expected version %d but got %d", 2, base.CurrentVersion())
	}

	events := base.GetUncommittedEvents()

	if len(events) != 2 {
		t.Fatalf("expected %d events but got %d", 2, len(events))
	}

	base.ClearUncommittedEvents()

	if base.Version() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.Version())
	}

	if base.CurrentVersion() != 0 {
		t.Errorf("expected version %d but got %d", 0, base.CurrentVersion())
	}

	events = base.GetUncommittedEvents()

	if len(events) != 0 {
		t.Fatalf("expected no events but got %d", len(events))
	}
}
