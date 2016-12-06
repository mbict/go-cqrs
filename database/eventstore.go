package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"strings"
)

var ErrNoAffectedRows = errors.New("No affected rows")

type EventStore struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	selectStmt *sql.Stmt
}

func NewDatabaseEventStore(db *sql.DB) cqrs.EventStore {

	//insertStmt, err := db.Prepare("INSERT INTO events (aggregate_id, aggregate_type, type, data, version, created) VALUES ($1, $2, $3, $4, $5, NOW())")
	insertStmt, err := db.Prepare("INSERT INTO events (aggregate_id, aggregate_type, type, data, version, created) VALUES (?,?,?,?,?, NOW())")
	if err != nil {
		log.Fatal(err)
	}

	//selectStmt, err := db.Prepare("SELECT aggregate_id, type, data, version, created FROM events WHERE aggregate_id = $1 AND aggregate_type = $2 ORDER BY version ASC")
	selectStmt, err := db.Prepare("SELECT aggregate_id, type, data, version, created FROM events WHERE aggregate_id = ? AND aggregate_type = ? ORDER BY version ASC")
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

func (s *EventStore) FindStream(aggregateTypes []string, aggregateIds []uuid.UUID, eventTypes []string) (cqrs.EventStream, error) {

	bindVars := []interface{}{}
	wheres := []string{}

	if l := len(aggregateIds); l >= 1 {

		//wheres = append(wheres, "aggregate_id IN ("+createEnumeratedBindParams(1, l)+")")
		wheres = append(wheres, "aggregate_id IN (?"+strings.Repeat(", ?", len(aggregateIds)-1)+")")
		for _, v := range aggregateIds {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(aggregateTypes); l >= 1 {
		//wheres = append(wheres, "aggreagate_type IN ("+createEnumeratedBindParams(len(bindVars)+1, l)+")")
		wheres = append(wheres, "aggreagate_type IN (?"+strings.Repeat(", ?", len(aggregateTypes)-1)+")")
		for _, v := range aggregateTypes {
			bindVars = append(bindVars, v)
		}
	}

	if l := len(eventTypes); l >= 1 {
		//wheres = append(wheres, "type IN ("+createEnumeratedBindParams(len(bindVars)+1, l)+")")
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
	fmt.Println(query, rows, err, bindVars)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return newDatabaseEventStream(rows), nil
}

func (s *EventStore) LoadStream(aggregateType string, aggregateId uuid.UUID) (cqrs.EventStream, error) {

	rows, err := s.selectStmt.Query(aggregateId, aggregateType)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return newDatabaseEventStream(rows), nil
}

func (s *EventStore) WriteEvent(aggregateType string, event cqrs.Event) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	result, err := s.insertStmt.Exec(MysqlUUID( event.AggregateID()), aggregateType, event.EventType(), payload, event.Version())
	if err != nil {
		return err
	}

	if affectedRows, err := result.RowsAffected(); err != nil || affectedRows == 0 {
		return ErrNoAffectedRows
	}

	return nil
}
