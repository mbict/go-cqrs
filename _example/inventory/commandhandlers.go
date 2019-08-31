package main

import (
	"context"
	"errors"
	"github.com/mbict/go-commandbus"
)

var ErrInventoryNameNotUnique = errors.New("inventory name not unique")

// an way to handle validation checks outside the aggregate
// you could provide a domain service or like this example a repository
func UniqueInventoryItemNameCommandHandlerMiddleware(repository InventoryNameRepository) commandbus.CommandHandler {
	return commandbus.CommandHandlerFunc(func(ctx context.Context, command commandbus.Command) error {
		switch c := command.(type) {
		case CreateInventoryItem:
			if repository.FindByName(c.Name) != nil {
				return ErrInventoryNameNotUnique
			}

		case RenameInventoryItem:
			item := repository.FindByName(c.Name)
			if item != nil && item.Id != c.AggregateId().String() {
				return ErrInventoryNameNotUnique
			}
		}

		return nil
	})
}
