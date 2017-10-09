package cqrs

// Validate is the interface an aggregate command can implement to perform
// validation prior to executing domain logic
type Validate interface {
	Validate() error
}
