package projection

import (
	"errors"
	"github.com/mbict/go-cqrs"
	"testing/base/domain/event"
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
		&event.ItemCreated{},
		&event.ItemTitleUpdated{},
		&event.ItemPriceUpdated{},
		&event.ItemDeleted{},
	}
}

//HandleEvent will apply the event
func (p *ItemsProjection) HandleEvent(event cqrs.Event) error {
	switch e := event.(type) {
	case *event.ItemCreated:
		return p.handleItemCreated(e)

	case *event.ItemTitleUpdated:
		return p.handleItemTitleUpdated(e)

	case *event.ItemPriceUpdated:
		return p.handleItemPriceUpdated(e)

	case *event.ItemDeleted:
		return p.handleItemDeleted(e)

	}
	return cqrs.ErrUnknownEvent
}


func (p *ItemsProjection) handleItemCreated(event *event.ItemCreated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemTitleUpdated(event *event.ItemTitleUpdated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemPriceUpdated(event *event.ItemPriceUpdated) error{
	//todo: implement event handling for this projection
	return nil
}

func (p *ItemsProjection) handleItemDeleted(event *event.ItemDeleted) error{
	//todo: implement event handling for this projection
	return nil
}
