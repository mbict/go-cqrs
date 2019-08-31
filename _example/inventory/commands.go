package main

import (
	"github.com/mbict/go-cqrs/v4"
)

type DeactivateInventoryItem struct {
	InventoryItemId cqrs.AggregateId
}

func (DeactivateInventoryItem) CommandName() string {
	return "deactivate_inventory_item"
}

func (c DeactivateInventoryItem) AggregateId() cqrs.AggregateId {
	return cqrs.NewAggregateId(c.InventoryItemId)
}

type CreateInventoryItem struct {
	InventoryItemId cqrs.AggregateId
	Name            string
}

func (CreateInventoryItem) CommandName() string {
	return "create_inventory_item"
}

func (c CreateInventoryItem) AggregateId() cqrs.AggregateId {
	return cqrs.NewAggregateId(c.InventoryItemId)
}

type RenameInventoryItem struct {
	InventoryItemId cqrs.AggregateId
	Name            string
}

func (RenameInventoryItem) CommandName() string {
	return "rename_inventory_item"
}

func (c RenameInventoryItem) AggregateId() cqrs.AggregateId {
	return cqrs.NewAggregateId(c.InventoryItemId)
}

type CheckInItemsToInventory struct {
	InventoryItemId cqrs.AggregateId
	Count           int
}

func (CheckInItemsToInventory) CommandName() string {
	return "checkin_items_to_inventory"
}

func (c CheckInItemsToInventory) AggregateId() cqrs.AggregateId {
	return cqrs.NewAggregateId(c.InventoryItemId)
}

type RemoveItemsFromInventory struct {
	InventoryItemId cqrs.AggregateId
	Count           int
}

func (RemoveItemsFromInventory) CommandName() string {
	return "remove_items_from_inventory"
}

func (c RemoveItemsFromInventory) AggregateId() cqrs.AggregateId {
	return cqrs.NewAggregateId(c.InventoryItemId)
}
