package cqrs

import (
	"testing"
)

func TestAggregateContext_AggregateId(t *testing.T) {
	id := NewIntAggregateId(12)
	ctx := NewAggregateContext(id, 0)

	if ctx.AggregateId().String() != "12" {
		t.Errorf("expected aggregateId %s but got %s", id.String(), ctx.AggregateId().String())
	}
}

func TestAggregateContext_incrementVersion(t *testing.T) {
	id := NewIntAggregateId(12)
	ctx := NewAggregateContext(id, 123)

	if ctx.Version() != 123 {
		t.Errorf("expected version %d but got %d", 123, ctx.Version())
	}

	if ctx.OriginalVersion() != 123 {
		t.Errorf("expected original version %d but got %d", 123, ctx.OriginalVersion())
	}

	ctx.incrementVersion()

	if ctx.Version() != 124 {
		t.Errorf("expected version %d but got %d", 123, ctx.Version())
	}

	if ctx.OriginalVersion() != 124 {
		t.Errorf("expected original version %d but got %d", 123, ctx.OriginalVersion())
	}
}

func TestAggregateContext_OrignalVersionShoulReturnCommitedVersion(t *testing.T) {
	id := NewIntAggregateId(12)
	ctx := NewAggregateContext(id, 123)

	if ctx.Version() != 123 {
		t.Errorf("expected version %d but got %d", 123, ctx.Version())
	}

	if ctx.OriginalVersion() != 123 {
		t.Errorf("expected original version %d but got %d", 123, ctx.OriginalVersion())
	}

	ctx.StoreEvent(eventA{})

	if ctx.OriginalVersion() != 123 {
		t.Errorf("expected original version %d but got %d", 123, ctx.OriginalVersion())
	}

	if ctx.Version() != 124 {
		t.Errorf("expected version %d but got %d", 124, ctx.Version())
	}
}

func TestAggregateContext_EventsHandling(t *testing.T) {
	id := NewIntAggregateId(12)
	ctx := NewAggregateContext(id, 0)

	ctx.StoreEvent(&eventA{})
	ctx.StoreEvent(&eventB{})

	events := ctx.getUncommittedEvents()

	if len(events) != 2 {
		t.Fatalf("expected %d events but got %d", 2, len(events))
	}

	ctx.clearUncommittedEvents()
	events = ctx.getUncommittedEvents()

	if len(events) != 0 {
		t.Fatalf("expected no events but got %d", len(events))
	}
}
