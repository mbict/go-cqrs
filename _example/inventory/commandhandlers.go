package main

import (
	"context"
	"errors"
	"github.com/mbict/go-commandbus"
	"github.com/satori/go.uuid"
)

var ErrInventoryNAmeNotUnique = errors.New("inventory name not unique")

// an way to handle validation checks outside the aggregate
// you could provide a domain service or like this example a repository
func UniqueInventoryItemNameCommandHandlerMiddleware(repository InventoryNameRepository) commandbus.CommandHandler {
	return commandbus.CommandHandlerFunc(func(ctx context.Context, command commandbus.Command) error {
		switch c := command.(type) {
		case CreateInventoryItem:
			if repository.FindByName(c.Name) != nil {
				return ErrInventoryNAmeNotUnique
			}

		case RenameInventoryItem:
			item := repository.FindByName(c.Name)
			if item != nil && !uuid.Equal(item.Id, c.AggregateId()) {
				return ErrInventoryNAmeNotUnique
			}
		}

		return nil
	})
}
