package projections

import (
	"errors"
	"github.com/mbict/go-cqrs"
	"testing/base/events"
)

type Items struct {

}

type ItemsProjection struct {}

func NewItemsProjection() *ItemsProjection {
	return &ItemsProjection{}
}

//HandlesEvents returns the events handled by this projection
func (p *ItemsProjection) HandlesEvents() []cqrs.Event {
	return []cqrs.Event{
		&events.ItemCreated{},
		&events.ItemTitleUpdated{},
		&events.ItemPriceUpdated{},
		&events.ItemDeleted{},
	}
}

//HandleEvent will apply the event
func (p *ItemsProjection) HandleEvent(event cqrs.Event) error {
	switch e := event.(type) {
	case *events.ItemCreated:
		return p.handleItemCreated(e)

	case *events.ItemTitleUpdated:
		return p.handleItemTitleUpdated(e)

	case *events.ItemPriceUpdated:
		return p.handleItemPriceUpdated(e)

	case *events.ItemDeleted:
		return p.handleItemDeleted(e)

	}
	return cqrs.ErrUnknownEvent
}


func (p *ItemsProjection) handleItemCreated(event *events.ItemCreated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemTitleUpdated(event *events.ItemTitleUpdated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemPriceUpdated(event *events.ItemPriceUpdated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemDeleted(event *events.ItemDeleted) error{
	//todo: implement event handling for this projection
	return nil
}
