package main

import (
	"github.com/mbict/go-cqrs/v4"
)

const (
	ItemDeactivated cqrs.EventType = "inventory:item_deactivated"
	ItemCreated     cqrs.EventType = "inventory:item_created"
	ItemRenamed     cqrs.EventType = "inventory:item_renamed"
	ItemsCheckedIn  cqrs.EventType = "inventory:items_checked_in"
	ItemsRemoved    cqrs.EventType = "inventory:items_removed"
)

type InventoryItemDeactivated struct {
}

func (InventoryItemDeactivated) EventType() cqrs.EventType {
	return ItemDeactivated
}

type InventoryItemCreated struct {
	Name string `json:"name"`
}

func (InventoryItemCreated) EventType() cqrs.EventType {
	return ItemCreated
}

type InventoryItemRenamed struct {
	NewName string `json:"new_name"`
}

func (InventoryItemRenamed) EventType() cqrs.EventType {
	return ItemRenamed
}

type ItemsCheckedInToInventory struct {
	Count int `json:"Count"`
}

func (ItemsCheckedInToInventory) EventType() cqrs.EventType {
	return ItemsCheckedIn
}

type ItemsRemovedFromInventory struct {
	Count int `json:"Count"`
}

func (ItemsRemovedFromInventory) EventType() cqrs.EventType {
	return ItemsRemoved
}
