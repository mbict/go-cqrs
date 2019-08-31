package main

// read models
type InventoryName struct {
	Id   Id
	Name string
}

type InventoryItem struct {
	Id    Id
	Name  string
	Count int
}

// interfaces of the repositories
type InventoryNameRepository interface {
	FindById(Id) *InventoryName
	FindByName(string) *InventoryName
	Save(*InventoryName) error
	Delete(Id) error
}

type InventoryItemRepository interface {
	FindById(Id) *InventoryItem
	Save(*InventoryItem) error
	Delete(Id) error
}

//actual naive implementation of the interface of the InventoryNameRepository
type inventoryNameRepository struct {
	names map[string]Id
	ids   map[Id]string
}

func (r *inventoryNameRepository) FindById(id Id) *InventoryName {
	if name, ok := r.ids[id]; ok {
		return &InventoryName{
			Id:   id,
			Name: name,
		}
	}
	return nil
}

func (r *inventoryNameRepository) FindByName(name string) *InventoryName {
	if id, ok := r.names[name]; ok {
		return &InventoryName{
			Id:   id,
			Name: name,
		}
	}
	return nil
}

func (r *inventoryNameRepository) Save(item *InventoryName) error {
	if name, ok := r.ids[item.Id]; ok {
		delete(r.names, name)
	}
	r.ids[item.Id] = item.Name
	r.names[item.Name] = item.Id
	return nil
}

func (r *inventoryNameRepository) Delete(id Id) error {
	if name, ok := r.ids[id]; ok {
		delete(r.names, name)
		delete(r.ids, id)
	}
	return nil
}

//actual naive implementation of the interface for InventoryItemRepository
type inventoryItemRepository struct {
	items map[Id]*InventoryItem
}

func (r *inventoryItemRepository) FindById(id Id) *InventoryItem {
	if item, ok := r.items[id]; ok {
		//we copy the item
		return &InventoryItem{
			Id:    item.Id,
			Name:  item.Name,
			Count: item.Count,
		}
	}
	return nil
}

func (r *inventoryItemRepository) Save(item *InventoryItem) error {
	r.items[item.Id] = item
	return nil
}

func (r *inventoryItemRepository) Delete(id Id) error {
	delete(r.items, id)
	return nil
}

// repository constructors
func NewInventoryNameRepository() InventoryNameRepository {
	return &inventoryNameRepository{
		names: make(map[string]Id),
		ids:   make(map[Id]string),
	}
}

func NewInventoryItemRepository() InventoryItemRepository {
	return &inventoryItemRepository{
		items: make(map[Id]*InventoryItem),
	}
}
