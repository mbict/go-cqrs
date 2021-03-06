package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mbict/go-cqrs/v4"
	"github.com/mbict/go-cqrs/v4/database"
	"log"
	"strings"
)

var ErrNoAffectedRows = errors.New("No affected rows")

type EventStore struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	selectStmt *sql.Stmt
}

func NewDatabaseEventStore(db *sql.DB) cqrs.EventStore {

	insertStmt, err := db.Prepare("INSERT INTO events (aggregate_id, aggregate_type, type, data, version, created) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Fatal(err)
	}

	selectStmt, err := db.Prepare("SELECT aggregate_id, type, data, version, created FROM events WHERE aggregate_id = $1 AND aggregate_type = $2 AND version > $3 ORDER BY version ASC")
	if err != nil {
		log.Fatal(err)
	}

	return &EventStore{
		db:         db,
		insertStmt: insertStmt,
		selectStmt: selectStmt,
	}
}

func createEnumeratedBindParams(offset, length int) string {
	result := fmt.Sprintf("$%d", offset)
	for i := offset + 1; i < offset+length; i++ {
		result = fmt.Sprintf("%s, $%d", result, i)
	}
	return result
}

func (s *EventStore) FindStream(aggregateTypes []string, aggregateIds []cqrs.AggregateId, eventTypes []string) (cqrs.EventStream, error) {

	bindVars := []interface{}{}
	wheres := []string{}

	if l := len(aggregateIds); l >= 1 {

		wheres = append(wheres, "aggregate_id IN ("+createEnumeratedBindParams(1, l)+")")
		for _, v := range aggregateIds {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(aggregateTypes); l >= 1 {
		wheres = append(wheres, "aggreagate_type IN ("+createEnumeratedBindParams(len(bindVars)+1, l)+")")
		for _, v := range aggregateTypes {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(eventTypes); l >= 1 {
		wheres = append(wheres, "type IN ("+createEnumeratedBindParams(len(bindVars)+1, l)+")")
		for _, v := range eventTypes {
			bindVars = append(bindVars, v)
		}
	}

	query := "SELECT aggregate_id, type, data, version, created FROM events"
	conditions := strings.Join(wheres, " AND ")
	if conditions != "" {
		query = query + " WHERE " + conditions
	}

	rows, err := s.db.Query(query+" ORDER BY id", bindVars...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return database.NewDatabaseEventStream(rows), nil
}

func (s *EventStore) LoadStream(aggregateType string, aggregateId cqrs.AggregateId, version int) (cqrs.EventStream, error) {

	rows, err := s.selectStmt.Query(aggregateId, aggregateType, version)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return database.NewDatabaseEventStream(rows), nil
}

func (s *EventStore) WriteEvent(aggregateType string, events ...cqrs.Event) error {
	tx, err := s.db.Begin()
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	for _, event := range events {
		payload, err := json.Marshal(event.Data())
		if err != nil {
			return err
		}

		result, err := tx.Stmt(s.insertStmt).Exec(event.AggregateId(), aggregateType, event.EventType(), payload, event.Version(), event.Timestamp())
		if err != nil {
			return err
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affectedRows == 0 {
			return ErrNoAffectedRows
		}
	}

	tx.Commit()
	return nil
}
