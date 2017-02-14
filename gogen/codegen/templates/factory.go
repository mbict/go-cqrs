package domain

import (
	"database/sql"
	"github.com/mbict/go-cqrs"
	"github.com/mbict/go-cqrs/database"
	"github.com/satori/go.uuid"
	"stalling/api/domain/aggregates"
	"stalling/api/domain/events"
	"stalling/api/domain/projections"
	"stalling/api/domain/validators"
	"stalling/api/repository"
)

type CQRSFactory interface {
	EventFactory() cqrs.EventFactory
	AggregateFactory() cqrs.AggregateFactory
}

type domainFactory struct {
	aggregateFactory cqrs.AggregateFactory
	eventFactory     cqrs.EventFactory
}

func NewCQRSFactory(db *sql.DB) (CQRSFactory, error) {

	aggregateFactory := cqrs.NewCallbackAggregateFactory()
	eventFactory := cqrs.NewCallbackEventFactory()

	//init aggregate and event factories
	initStallingFactories(db, aggregateFactory, eventFactory)
	initParkingFactories(db, aggregateFactory, eventFactory)
	initItemFactories(db, aggregateFactory, eventFactory)
	initCustomerFactories(db, aggregateFactory, eventFactory)

	return &domainFactory{
		aggregateFactory: aggregateFactory,
		eventFactory:     eventFactory,
	}, nil
}

func (f *domainFactory) EventFactory() cqrs.EventFactory {
	return f.eventFactory
}

func (f *domainFactory) AggregateFactory() cqrs.AggregateFactory {
	return f.aggregateFactory
}

func initStallingFactories(db *sql.DB, aggregateFactory *cqrs.CallbackAggregateFactory, eventFactory *cqrs.CallbackEventFactory) {

	stallingRepository := repository.NewStallingRepository(db)
	stallingValidator := validators.NewStallingDomainValidator(stallingRepository)
	aggregateFactory.RegisterCallback(func(id uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewStallingAggregate(id, stallingValidator)
	})

	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingCreated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingDescriptionChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingDomainChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingNameChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingDeleted{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingLocationAdded{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingLocationRenamed{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.StallingLocationRemoved{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
}

func initItemFactories(db *sql.DB, aggregateFactory *cqrs.CallbackAggregateFactory, eventFactory *cqrs.CallbackEventFactory) {

	aggregateFactory.RegisterCallback(func(id uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewItemAggregate(id)
	})

	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemCreated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemNameChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemCodeChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemHasCustomerAssigned{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemBookedOut{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemBookedIn{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemBookingTransfered{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemPickupDateAdded{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemPickupDateRemoved{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ItemDeleted{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
}

func initParkingFactories(db *sql.DB, aggregateFactory *cqrs.CallbackAggregateFactory, eventFactory *cqrs.CallbackEventFactory) {

	aggregateFactory.RegisterCallback(func(id uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewParkingAggregate(id)
	})

	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingCreated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingNameChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingCodeChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingLotAdded{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingLotRenamed{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingLotRemoved{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.ParkingDeleted{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
}

func initCustomerFactories(db *sql.DB, aggregateFactory *cqrs.CallbackAggregateFactory, eventFactory *cqrs.CallbackEventFactory) {

	aggregateFactory.RegisterCallback(func(id uuid.UUID) cqrs.AggregateRoot {
		return aggregates.NewCustomerAggregate(id)
	})

	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.CustomerCreated{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.CustomerNameChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.CustomerNumberChanged{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.CustomerHasLoginUserAssigned{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
	eventFactory.RegisterCallback(func(aggregateId uuid.UUID, version int) cqrs.Event {
		return &events.CustomerDeleted{
			EventBase: cqrs.NewEventBase(aggregateId, version),
		}
	})
}

//read models (projections)
func (d *domain) withProjections(db *sql.DB) {

	StallingProjection := projections.NewStallingProjection(db)
	d.eventBus.AddHandler(StallingProjection, StallingProjection.HandlesEvents()...)

	LocationProjection := projections.NewLocationProjection(db)
	d.eventBus.AddHandler(LocationProjection, LocationProjection.HandlesEvents()...)

	BarcodeProjection := projections.NewCodeProjection(db)
	d.eventBus.AddHandler(BarcodeProjection, BarcodeProjection.HandlesEvents()...)

	ParkingProjection := projections.NewParkingsProjection(db)
	d.eventBus.AddHandler(ParkingProjection, ParkingProjection.HandlesEvents()...)

	CustomerProjection := projections.NewCustomersProjection(db)
	d.eventBus.AddHandler(CustomerProjection, CustomerProjection.HandlesEvents()...)

	ItemProjection := projections.NewItemsProjection(db)
	d.eventBus.AddHandler(ItemProjection, ItemProjection.HandlesEvents()...)

	PickupProjection := projections.NewPickupsProjection(db)
	d.eventBus.AddHandler(PickupProjection, PickupProjection.HandlesEvents()...)
}


