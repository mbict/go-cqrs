package domain

import (
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"testing/base/event"
)

func RegisterEventFactory(factory *cqrs.CallbackEventFactory) error {
	itemDeletedEventFactory := func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &event.ItemDeleted{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	}
	if err := factory.RegisterCallback(itemDeletedEventFactory); err != nil {
		return err
	}

	itemPriceUpdatedEventFactory := func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &event.ItemPriceUpdated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	}
	if err := factory.RegisterCallback(itemPriceUpdatedEventFactory); err != nil {
		return err
	}

	itemTitleUpdatedEventFactory := func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &event.ItemTitleUpdated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	}
	if err := factory.RegisterCallback(itemTitleUpdatedEventFactory); err != nil {
		return err
	}

	test2CreatedEventFactory := func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &event.Test2Created{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	}
	if err := factory.RegisterCallback(test2CreatedEventFactory); err != nil {
		return err
	}

	return nil
}
