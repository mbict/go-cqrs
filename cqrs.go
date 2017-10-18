package cqrs

import (
	"errors"
)

var ErrUnknownCommand = errors.New("Cannot handle unknown Command")
var ErrUnknownEvent = errors.New("Cannot handle unknown Event")
