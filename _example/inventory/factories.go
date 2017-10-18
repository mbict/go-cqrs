package main

import (
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
)

//factory for the aggregate
func inventoryItemAggregateFactory(id uuid.UUID) cqrs.Aggregate {
	return &InventoryItemAggregate{
		AggregateBase: cqrs.NewAggregateBase(id),
	}
}

//factories for the events
func inventoryItemDeactivatedFactory(id uuid.UUID, version int) cqrs.Event {
	return &InventoryItemDeactivated{
		EventBase: cqrs.NewEventBase(id, version),
	}
}
func inventoryItemCreatedFactory(id uuid.UUID, version int) cqrs.Event {
	return &InventoryItemCreated{
		EventBase: cqrs.NewEventBase(id, version),
	}
}
func inventoryItemRenamedFactory(id uuid.UUID, version int) cqrs.Event {
	return &InventoryItemRenamed{
		EventBase: cqrs.NewEventBase(id, version),
	}
}
func itemsCheckedInToInventoryFactory(id uuid.UUID, version int) cqrs.Event {
	return &ItemsCheckedInToInventory{
		EventBase: cqrs.NewEventBase(id, version),
	}
}
func itemsRemovedFromInventoryFactory(id uuid.UUID, version int) cqrs.Event {
	return &ItemsRemovedFromInventory{
		EventBase: cqrs.NewEventBase(id, version),
	}
}

//event factory
var eventFactory = cqrs.NewCallbackEventFactory()

func init() {
	mustSucceed(eventFactory.RegisterCallback(inventoryItemDeactivatedFactory))
	mustSucceed(eventFactory.RegisterCallback(inventoryItemCreatedFactory))
	mustSucceed(eventFactory.RegisterCallback(inventoryItemRenamedFactory))
	mustSucceed(eventFactory.RegisterCallback(itemsCheckedInToInventoryFactory))
	mustSucceed(eventFactory.RegisterCallback(itemsRemovedFromInventoryFactory))
}

func mustSucceed(e error) {
	if e != nil {
		panic(e)
	}
}
