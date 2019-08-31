package main

import (
	"fmt"
	"github.com/mbict/go-commandbus"
	"github.com/mbict/go-cqrs/v4"
	"github.com/mbict/go-cqrs/v4/memory"
	"github.com/mbict/go-eventbus"
	"github.com/satori/go.uuid"
)

var (
	//commandbus for handling the commands
	commandBus commandbus.CommandBus

	//read repository used by application, projections and command middleware
	itemRepository     InventoryItemRepository
	itemNameRepository InventoryNameRepository
)

//bootstrapping
func init() {
	commandBus = commandbus.New()
	itemRepository = NewInventoryItemRepository()
	itemNameRepository = NewInventoryNameRepository()

	//in memory eventbus for dispatching events
	eventBus := eventbus.New()

	// hook to eventbus when successful store of event to the eventstore we
	// will publish the event to the in memory eventbus.
	eventPublisher := cqrs.EventPublisherFunc(func(events ...cqrs.Event) error {
		for _, event := range events {
			if err := eventBus.Publish(event.(eventbus.Event)); err != nil {
				return err
			}
		}
		return nil
	})

	// initilaze inmemory eventbus and aggregateRepository for the InventoryItemAggregate
	var eventStore cqrs.EventStore
	{
		eventStore = memory.NewMemoryEventStore()
		eventStore = cqrs.NewEventPublishingEventStore(eventPublisher, eventStore)
	}

	//aggregateRepository := cqrs.NewAggregateRepository(eventStore, cqrs.DefaultAggregateBuilder(inventoryItemAggregateFactory), eventFactory, eventbusNotifyHook)

	//snapshot version of the repository
	snapshotStore := memory.NewSnapshotStore()
	aggregateRepository := cqrs.NewSnapshotAggregateRepository(
		eventStore,
		snapshotStore,
		2, //create a snapshot everytime we differ 2 events or more
		cqrs.SnapshotAggregateBuilder(inventoryItemAggregateFactory, snapshotStore),
		eventFactory)

	// aggregate command handler is the command handler who is responsible for
	// - creating the aggregate
	// - load all the events from the eventstore into the aggregate
	// - pass the command to the aggregate command handling function
	// - store the generated events upon successful handling of the command
	aggregateCommandHandler := cqrs.AggregateCommandHandler(aggregateRepository)

	// extra domain logic for checking the uniqueness of a inventory item name
	// is handled by a command middleware, or actual chained command handlers who will only execute
	// if this middleware succeeds execution.
	// We will add this handler for the command CreateInventoryItem and RenameInventoryItem.
	uniqueNameMiddleware := UniqueInventoryItemNameCommandHandlerMiddleware(itemNameRepository)

	//register our commands
	commandBus.Register(CreateInventoryItem{}, commandbus.ChainHandler(aggregateCommandHandler, uniqueNameMiddleware))
	commandBus.Register(DeactivateInventoryItem{}, aggregateCommandHandler)
	commandBus.Register(RenameInventoryItem{}, commandbus.ChainHandler(aggregateCommandHandler, uniqueNameMiddleware))
	commandBus.Register(CheckInItemsToInventory{}, aggregateCommandHandler)
	commandBus.Register(RemoveItemsFromInventory{}, aggregateCommandHandler)

	// example for the projection part of the read models
	// this is now directly chained to the internal memory bus, but this could be connected
	// to a redis pubsub, kafka or an other messaging service.
	uniqueInventoryNamesProjector := NewInventoryNameProjector(itemNameRepository)
	inventoryItemProjector := NewInventoryProjector(itemRepository)

	// subscribe to the in memory eventbus, to all events
	eventBus.Subscribe(cqrs.EventbusWrapper(inventoryItemProjector))

	// subscribe to the in memory eventbus to only these specific events
	eventBus.Subscribe(cqrs.EventbusWrapper(uniqueInventoryNamesProjector), uniqueInventoryNamesProjector.HandlesEvent()...)
}

//rough example
func main() {
	// the unique id of the first inventoryItem
	idFirstItem := uuid.Must(uuid.NewV4())

	//we want to create a new inventory item
	commandCreate := CreateInventoryItem{
		InventoryItemId: idFirstItem,
		Name:            "battery",
	}

	//lets feed it to the commandbus
	err := commandBus.Handle(nil, commandCreate)
	if err != nil {
		panic(err)
	}

	//check if the projection has inserted the item into the read repository
	item := itemRepository.FindById(idFirstItem)
	if item == nil {
		panic("aahhww something went wrong")
	}
	fmt.Printf("[SUCCESS] we got a inventory item from the repo:\n%#v\n", item)

	//oops did i just make a typo in the name
	commandRenameItem := RenameInventoryItem{
		InventoryItemId: idFirstItem,
		Name:            "duracell battery",
	}

	//it went so great the first time lets try if we are se lucky the second time
	err = commandBus.Handle(nil, commandRenameItem)
	if err != nil {
		panic(err)
	}

	//check if the projection has updated our aggregates name
	item = itemRepository.FindById(idFirstItem)
	if item.Name != "duracell battery" {
		panic("the projection was lazy and did not update the inventory items name")
	}
	fmt.Printf("[SUCCESS] inventory item item after renaming:\n%#v\n", item)

	//Lets checkin items to the inventory
	commandCheckIn := CheckInItemsToInventory{
		InventoryItemId: idFirstItem,
		Count:           500,
	}

	err = commandBus.Handle(nil, commandCheckIn)
	if err != nil {
		panic(err)
	}

	//check if the projection has updated our aggregates name
	item = itemRepository.FindById(idFirstItem)
	if item.Count != 500 {
		panic("there should be 500 checked in items in the inventory, but there was a thief who took some")
	}
	fmt.Printf("[SUCCESS] inventory item item after checking in items:\n%#v\n", item)

	//now lets test the middleware if the unique names work
	newId := uuid.Must(uuid.NewV4())
	createCommandWithDuplicateName := CreateInventoryItem{
		InventoryItemId: newId,
		Name:            "duracell battery",
	}

	err = commandBus.Handle(nil, createCommandWithDuplicateName)
	if err == nil {
		panic("we expected a error here about the uniqueness of the item name")
	}
	fmt.Printf("[SUCCESS] we got a error that the command is not processed: `%s` and that is just what we want\n", err)

	//and now add an remove some item to the inventory
	commandBus.Handle(nil, CheckInItemsToInventory{
		InventoryItemId: idFirstItem,
		Count:           657,
	})
	commandBus.Handle(nil, RemoveItemsFromInventory{
		InventoryItemId: idFirstItem,
		Count:           123,
	})

	//and the final result is 500 + 657 - 123 = 1034
	item = itemRepository.FindById(idFirstItem)
	if item.Count != 1034 {
		panic("we expected the outcome to be 1034 but its something different")
	}
	fmt.Printf("[SUCCESS] the inventory is changed in a good way, no thiefs today:\n%#v\n", item)

}
