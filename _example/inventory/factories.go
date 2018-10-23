package main

import (
	"github.com/mbict/go-cqrs"
)

//factory for the aggregate
func inventoryItemAggregateFactory(context cqrs.AggregateContext) cqrs.Aggregate {
	return &InventoryItemAggregate{
		AggregateContext: context,
	}
}

//factories for the events
func inventoryItemDeactivatedFactory() cqrs.EventData {
	return &InventoryItemDeactivated{}
}

func inventoryItemCreatedFactory() cqrs.EventData {
	return &InventoryItemCreated{}
}

func inventoryItemRenamedFactory() cqrs.EventData {
	return &InventoryItemRenamed{}
}

func itemsCheckedInToInventoryFactory() cqrs.EventData {
	return &ItemsCheckedInToInventory{}
}

func itemsRemovedFromInventoryFactory() cqrs.EventData {
	return &ItemsRemovedFromInventory{}
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
