package main

import (
	"github.com/mbict/go-cqrs/v4"
)

//InventoryNameProjector is used for checking the uniqueness of the inventory name in
//the command middleware
type InventoryNameProjector struct {
	repository InventoryNameRepository
}

func (p *InventoryNameProjector) HandlesEvent() []cqrs.EventType {
	return []cqrs.EventType{
		ItemCreated,
		ItemRenamed,
		ItemDeactivated,
	}
}

func (p *InventoryNameProjector) Handle(event cqrs.Event) error {
	switch data := event.Data().(type) {
	case InventoryItemCreated:
		item := &InventoryName{
			Id:   event.AggregateId(),
			Name: data.Name,
		}
		p.repository.Save(item)

	case InventoryItemRenamed:
		if item := p.repository.FindById(event.AggregateId()); item != nil {
			item.Name = data.NewName
			p.repository.Save(item)
		}

	case InventoryItemDeactivated:
		p.repository.Delete(event.AggregateId())
	}
	return nil
}

func NewInventoryNameProjector(repository InventoryNameRepository) *InventoryNameProjector {
	return &InventoryNameProjector{
		repository: repository,
	}
}

type InventoryProjector struct {
	repository InventoryItemRepository
}

func (p *InventoryProjector) Handle(event cqrs.Event) error {
	switch data := event.Data().(type) {
	case InventoryItemCreated:
		item := &InventoryItem{
			Id:   event.AggregateId(),
			Name: data.Name,
		}
		p.repository.Save(item)

	case InventoryItemRenamed:
		if item := p.repository.FindById(event.AggregateId()); item != nil {
			item.Name = data.NewName
			p.repository.Save(item)
		}

	case InventoryItemDeactivated:
		p.repository.Delete(event.AggregateId())

	case ItemsCheckedInToInventory:
		if item := p.repository.FindById(event.AggregateId()); item != nil {
			item.Count += data.Count
			p.repository.Save(item)
		}

	case ItemsRemovedFromInventory:
		if item := p.repository.FindById(event.AggregateId()); item != nil {
			item.Count -= data.Count
			p.repository.Save(item)
		}
	}
	return nil
}

func NewInventoryProjector(repository InventoryItemRepository) *InventoryProjector {
	return &InventoryProjector{
		repository: repository,
	}
}
