package main

import (
	"github.com/mbict/go-cqrs"
)

type InventoryItemDeactivated struct {
	cqrs.EventBase
}

func (InventoryItemDeactivated) EventName() string {
	return "inventory_item_deactivated"
}

type InventoryItemCreated struct {
	cqrs.EventBase
	Name string `json:"name"`
}

func (InventoryItemCreated) EventName() string {
	return "inventory_item_created"
}

type InventoryItemRenamed struct {
	cqrs.EventBase
	NewName string `json:"new_name"`
}

func (InventoryItemRenamed) EventName() string {
	return "inventory_item_renamed"
}

type ItemsCheckedInToInventory struct {
	cqrs.EventBase
	Count int `json:"count"`
}

func (ItemsCheckedInToInventory) EventName() string {
	return "items_checked_in_to_inventory"
}

type ItemsRemovedFromInventory struct {
	cqrs.EventBase
	Count int `json:"count"`
}

func (ItemsRemovedFromInventory) EventName() string {
	return "items_removed_from_inventory"
}
