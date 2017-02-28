package aggregate

import (
	"errors"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"testing/base/domain/command"
	"testing/base/domain/event"
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
	case *event.ItemCreated:
		a.applyItemCreated(e)

	case *event.ItemTitleUpdated:
		a.applyItemTitleUpdated(e)

	case *event.ItemPriceUpdated:
		a.applyItemPriceUpdated(e)

	case *event.ItemDeleted:
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
	case *command.CreateItem:
		return a.handleCreateItem(c)

	case *command.UpdateItem:
		return a.handleUpdateItem(c)

	case *command.DeleteItem:
		return a.handleDeleteItem(c)

	}
	return cqrs.ErrUnknownCommand
}

func (a *ItemAggregate) handleCreateItem(command *command.CreateItem) error {
	//todo: implement command handling/validation here
	return nil
}

func (a *ItemAggregate) handleUpdateItem(command *command.UpdateItem) error {
	//todo: implement command handling/validation here
	return nil
}

func (a *ItemAggregate) handleDeleteItem(command *command.DeleteItem) error {
	//todo: implement command handling/validation here
	return nil
}
