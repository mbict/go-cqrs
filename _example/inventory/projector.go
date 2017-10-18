package main

import "github.com/mbict/go-eventbus"

//InventoryNameProjector is used for checking the uniqueness of the inventory name in
//the command middleware
type InventoryNameProjector struct {
	repository InventoryNameRepository
}

func (p *InventoryNameProjector) Handle(event eventbus.Event) {
	switch e := event.(type) {
	case InventoryItemCreated:
		item := &InventoryName{
			Id:   e.AggregateId(),
			Name: e.Name,
		}
		p.repository.Save(item)

	case InventoryItemRenamed:
		if item := p.repository.FindById(e.AggregateId()); item != nil {
			item.Name = e.NewName
			p.repository.Save(item)
		}

	case InventoryItemDeactivated:
		p.repository.Delete(e.AggregateId())
	}
}

func NewInventoryNameProjector(repository InventoryNameRepository) *InventoryNameProjector {
	return &InventoryNameProjector{
		repository: repository,
	}
}

type InventoryProjector struct {
	repository InventoryItemRepository
}

func (p *InventoryProjector) Handle(event eventbus.Event) {
	switch e := event.(type) {
	case InventoryItemCreated:
		item := &InventoryItem{
			Id:   e.AggregateId(),
			Name: e.Name,
		}
		p.repository.Save(item)

	case InventoryItemRenamed:
		if item := p.repository.FindById(e.AggregateId()); item != nil {
			item.Name = e.NewName
			p.repository.Save(item)
		}

	case InventoryItemDeactivated:
		p.repository.Delete(e.AggregateId())

	case ItemsCheckedInToInventory:
		if item := p.repository.FindById(e.AggregateId()); item != nil {
			item.Count += e.Count
			p.repository.Save(item)
		}

	case ItemsRemovedFromInventory:
		if item := p.repository.FindById(e.AggregateId()); item != nil {
			item.Count -= e.Count
			p.repository.Save(item)
		}
	}
}

func NewInventoryProjector(repository InventoryItemRepository) *InventoryProjector {
	return &InventoryProjector{
		repository: repository,
	}
}
