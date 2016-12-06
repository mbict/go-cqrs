package cqrs

type Validate interface {
	Validate() error
}
