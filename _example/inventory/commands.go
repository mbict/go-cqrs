package main

import "github.com/satori/go.uuid"

type DeactivateInventoryItem struct {
	InventoryItemId uuid.UUID
}

func (DeactivateInventoryItem) CommandName() string {
	return "deactivate_inventory_item"
}

func (c DeactivateInventoryItem) AggregateId() uuid.UUID {
	return c.InventoryItemId
}

type CreateInventoryItem struct {
	InventoryItemId uuid.UUID
	Name            string
}

func (CreateInventoryItem) CommandName() string {
	return "create_inventory_item"
}

func (c CreateInventoryItem) AggregateId() uuid.UUID {
	return c.InventoryItemId
}

type RenameInventoryItem struct {
	InventoryItemId uuid.UUID
	Name            string
}

func (RenameInventoryItem) CommandName() string {
	return "rename_inventory_item"
}

func (c RenameInventoryItem) AggregateId() uuid.UUID {
	return c.InventoryItemId
}

type CheckInItemsToInventory struct {
	InventoryItemId uuid.UUID
	Count           int
}

func (CheckInItemsToInventory) CommandName() string {
	return "checkin_items_to_inventory"
}

func (c CheckInItemsToInventory) AggregateId() uuid.UUID {
	return c.InventoryItemId
}

type RemoveItemsFromInventory struct {
	InventoryItemId uuid.UUID
	Count           int
}

func (RemoveItemsFromInventory) CommandName() string {
	return "remove_items_from_invertory"
}

func (c RemoveItemsFromInventory) AggregateId() uuid.UUID {
	return c.InventoryItemId
}
