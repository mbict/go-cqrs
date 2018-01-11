package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/mbict/go-cqrs"
	"github.com/mbict/go-cqrs/database"
	"github.com/satori/go.uuid"
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
	insertStmt, err := db.Prepare("INSERT INTO events (aggregate_id, aggregate_type, type, data, version, created) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	selectStmt, err := db.Prepare("SELECT aggregate_id, type, data, version, created FROM events WHERE aggregate_id = ? AND aggregate_type = ? version > ? ORDER BY version ASC")
	if err != nil {
		log.Fatal(err)
	}

	return &EventStore{
		db:         db,
		insertStmt: insertStmt,
		selectStmt: selectStmt,
	}
}

func (s *EventStore) FindStream(aggregateTypes []string, aggregateIds []uuid.UUID, eventTypes []string) (cqrs.EventStream, error) {

	bindVars := []interface{}{}
	wheres := []string{}

	if l := len(aggregateIds); l >= 1 {
		wheres = append(wheres, "aggregate_id IN (?"+strings.Repeat(", ?", len(aggregateIds)-1)+")")
		for _, v := range aggregateIds {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(aggregateTypes); l >= 1 {
		wheres = append(wheres, "aggreagate_type IN (?"+strings.Repeat(", ?", len(aggregateTypes)-1)+")")
		for _, v := range aggregateTypes {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(eventTypes); l >= 1 {
		wheres = append(wheres, "type IN (?"+strings.Repeat(", ?", len(eventTypes)-1)+")")
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

func (s *EventStore) LoadStream(aggregateType string, aggregateId uuid.UUID, version int) (cqrs.EventStream, error) {
	rows, err := s.selectStmt.Query(MysqlUUID(aggregateId), aggregateType, version)
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
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		result, err := tx.Stmt(s.insertStmt).Exec(MysqlUUID(event.AggregateId()), aggregateType, event.EventName(), payload, event.Version(), event.OccurredAt())
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
