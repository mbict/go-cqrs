package cqrs

import (
	"errors"
)

var ErrUnknownCommand = errors.New("Cannot handle unkown Command")
var ErrUnknownEvent = errors.New("Cannot handle unkown Event")
