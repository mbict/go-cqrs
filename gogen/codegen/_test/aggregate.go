package aggregate

import (
	"errors"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"testing/base/command"
	"testing/base/event"
)

var (
	ErrItemAlreadyExists = errors.New("Item already exists")
)

type ItemAggregate struct {
	*cqrs.AggregateBase
}

func NewItemAggregate(id uuid.UUID) *ItemAggregate {
	return &ItemAggregate{
		AggregateBase: cqrs.NewAggregateBase(id),
	}
}

func (a *ItemAggregate) AggregateType() string {
	return "Item"
}

func (a *ItemAggregate) Apply(event cqrs.Event) error {
	switch e := event.(type) {
	case *events.ItemCreated:
		a.applyItemCreated(e)

	case *events.ItemTitleUpdated:
		a.applyItemTitleUpdated(e)

	case *events.ItemPriceUpdated:
		a.applyItemPriceUpdated(e)

	case *events.ItemDeleted:
		a.applyItemDeleted(e)

	}
	return cqrs.ErrUnknownEvent
}

func (a *ItemAggregate) applyItemCreated(event *event.ItemCreated) {
	//todo: implement apply logic here
}

func (a *ItemAggregate) applyItemTitleUpdated(event *event.ItemTitleUpdated) {
	//todo: implement apply logic here
}

func (a *ItemAggregate) applyItemPriceUpdated(event *event.ItemPriceUpdated) {
	//todo: implement apply logic here
}

func (a *ItemAggregate) applyItemDeleted(event *event.ItemDeleted) {
	//todo: implement apply logic here
}

func (a *ItemAggregate) HandleCommand(command cqrs.Command) error {

	switch c := command.(type) {
	case *commands.CreateItem:
		return a.handleCreateItem(c)

	case *commands.UpdateItem:
		return a.handleUpdateItem(c)

	case *commands.DeleteItem:
		return a.handleDeleteItem(c)

	}
	return cqrs.ErrUnknownCommand
}

func (a *ItemAggregate) handleCreateItem(command *commands.CreateItem) error {
	//todo: implement command handling/validation here
	return nil
}

func (a *ItemAggregate) handleUpdateItem(command *commands.UpdateItem) error {
	//todo: implement command handling/validation here
	return nil
}

func (a *ItemAggregate) handleDeleteItem(command *commands.DeleteItem) error {
	//todo: implement command handling/validation here
	return nil
}
