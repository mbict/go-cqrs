# go-cqrs
CQRS/es implementation in go

### Aggregates
todo

### Commands
todo

### Events
todo

### Event store and streams
todo

### Factories
todo

### Aggregate repository
todo

### Projections
todo

### Middleware for command handler
todo

### Event publishing
todo

### Examples
Examples are in the `_example` directory.

**Inventory example**
This is the GO version of the C# Simple CQRS example from Gregory Young.


### Passing events by value not by pointer reference
Events are always passed by value, never passed by pointer reference. This is to ensure immutability of the events data.

- The Load method of the aggregate repository will always convert pointer events and pass them by value to the aggregates Apply function.
- The Save method of the aggregate repository will convert any pointer referenced events and pass them by value to the aggregate Apply method and to the publish event hooks (PublishEventFunc)  

> Always refer to the class name of an event in your aggregates and projections.
> Never to the pointer variant, it will probably never be picked up

### EventFactory
The event factory must always return a pointer to the newly created event.
This is of the event stream needs to scan/unmarshal the data into the event instance and can only do this for pointer instances.
Later on the newly created event will be passed by value to the aggregates and or projections.

### Todo
- Test the domain aggregate repository.
- More real world examples.
- See how it performs in my projects.
- Documentation documentation documentation.