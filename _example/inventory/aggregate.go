package main

import (
	"errors"
	"github.com/mbict/go-cqrs"
)

var ErrAlreadyCreated = errors.New("already created")
var ErrNotEnoughInventoryItems = errors.New("not enough inventory items available")
var ErrDeactivate = errors.New("cannot deactivated")

type InventoryItemAggregate struct {
	*cqrs.AggregateBase

	//aggregate state
	count     int
	activated bool
}

func (InventoryItemAggregate) AggregateName() string {
	return "inventory_item"
}

func (a *InventoryItemAggregate) HandleCommand(command cqrs.Command) error {
	switch c := command.(type) {
	case CreateInventoryItem:
		return a.create(c.Name)
	case DeactivateInventoryItem:
		return a.deactivate()
	case RenameInventoryItem:
		return a.rename(c.Name)
	case CheckInItemsToInventory:
		return a.checkinItems(c.Count)
	case RemoveItemsFromInventory:
		return a.removeItems(c.Count)
	}
	return nil
}

func (a *InventoryItemAggregate) Apply(event cqrs.Event) error {
	switch e := event.(type) {
	case InventoryItemCreated:
		a.activated = true
	case InventoryItemRenamed:
		//we do nothing here intentionally
		//there is no need to change the state of the aggregate for this event
	case InventoryItemDeactivated:
		a.activated = false
	case ItemsCheckedInToInventory:
		a.count += e.Count //add items to the state count
	case ItemsRemovedFromInventory:
		a.count -= e.Count //remove items from the state count
	default:
		//unkown event, we return an error (if anyone cares)
		return cqrs.ErrUnknownEvent
	}
	return nil
}

func (a *InventoryItemAggregate) create(name string) error {
	//check if not created, simple check on version
	if a.Version() != 0 {
		return ErrAlreadyCreated
	}

	event := InventoryItemCreated{
		EventBase: cqrs.NewEventBaseFromAggregate(a),
		Name:      name,
	}
	a.StoreEvent(event)

	return nil
}

func (a *InventoryItemAggregate) rename(name string) error {
	event := InventoryItemRenamed{
		EventBase: cqrs.NewEventBaseFromAggregate(a),
		NewName:   name,
	}
	a.StoreEvent(event)

	return nil
}

func (a *InventoryItemAggregate) deactivate() error {
	//pre check
	if a.activated == false {
		return ErrDeactivate
	}

	event := InventoryItemDeactivated{
		EventBase: cqrs.NewEventBaseFromAggregate(a),
	}
	a.StoreEvent(event)

	return nil
}

func (a *InventoryItemAggregate) checkinItems(count int) error {
	event := ItemsCheckedInToInventory{
		EventBase: cqrs.NewEventBaseFromAggregate(a),
		Count:     count,
	}
	a.StoreEvent(event)

	return nil
}

func (a *InventoryItemAggregate) removeItems(count int) error {
	if count > a.count {
		return ErrNotEnoughInventoryItems
	}

	event := ItemsRemovedFromInventory{
		EventBase: cqrs.NewEventBaseFromAggregate(a),
		Count:     count,
	}
	a.StoreEvent(event)

	return nil
}
