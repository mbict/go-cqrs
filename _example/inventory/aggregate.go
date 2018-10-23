package main

import (
	"errors"
	"github.com/mbict/go-cqrs"
)

var ErrAlreadyCreated = errors.New("already created")
var ErrNotEnoughInventoryItems = errors.New("not enough inventory items available")
var ErrDeactivate = errors.New("cannot deactivated")

type InventoryItemAggregate struct {
	// the context is used to store the generated events version and aggregate version
	// incase of a snapshot repository is used this field needs to be unexported or ignored by json
	cqrs.AggregateContext `json:"-"`

	//aggregate state, serializable for snapshots
	Count     int
	Activated bool
}

func (InventoryItemAggregate) AggregateName() string {
	return "inventory_item"
}

func (a *InventoryItemAggregate) HandleCommand(command cqrs.Command) error {
	//simple pre check if first command always must be from specific type
	if _, ok := command.(CreateInventoryItem); !ok && a.Version() == 0 {
		return cqrs.ErrorAggregateNotFound(a.AggregateId().String())
	}

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
	switch e := event.Data().(type) {
	case InventoryItemCreated:
		a.Activated = true
	case InventoryItemRenamed:
		//we do nothing here intentionally
		//there is no need to change the state of the aggregate for this event
	case InventoryItemDeactivated:
		a.Activated = false
	case ItemsCheckedInToInventory:
		a.Count += e.Count //add items to the state Count
	case ItemsRemovedFromInventory:
		a.Count -= e.Count //remove items from the state Count
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

	a.StoreEvent(InventoryItemCreated{
		Name: name,
	})

	return nil
}

func (a *InventoryItemAggregate) rename(name string) error {
	a.StoreEvent(InventoryItemRenamed{
		NewName: name,
	})

	return nil
}

func (a *InventoryItemAggregate) deactivate() error {
	//pre check
	if a.Activated == false {
		return ErrDeactivate
	}

	a.StoreEvent(InventoryItemDeactivated{})

	return nil
}

func (a *InventoryItemAggregate) checkinItems(count int) error {
	a.StoreEvent(ItemsCheckedInToInventory{
		Count: count,
	})

	return nil
}

func (a *InventoryItemAggregate) removeItems(count int) error {
	if count > a.Count {
		return ErrNotEnoughInventoryItems
	}

	a.StoreEvent(ItemsRemovedFromInventory{
		Count: count,
	})

	return nil
}
