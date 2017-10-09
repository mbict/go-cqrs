package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"time"
)

var (
	ErrNoEventData = errors.New("Cannot scan, no event data")
)

type EventStream struct {
	rows *sql.Rows

	err         error
	version     int
	aggregateId uuid.UUID
	eventType   string
	data        sql.RawBytes
	timestamp   time.Time
}

func newDatabaseEventStream(rows *sql.Rows) cqrs.EventStream {
	return &EventStream{
		rows: rows,
	}
}

func (s *EventStream) AggregateId() uuid.UUID {
	return s.aggregateId
}

func (s *EventStream) EventName() string {
	return s.eventType
}

func (s *EventStream) Version() int {
	return s.version
}

func (s *EventStream) Next() bool {
	if s.rows.Next() {
		if s.err = s.rows.Scan(&s.aggregateId, &s.eventType, &s.data, &s.version, &s.timestamp); s.err != nil {
			return false
		}
		return true
	}

	if s.err == nil {
		s.err = s.rows.Err()
	}
	s.aggregateId = uuid.Nil
	s.eventType = ""
	s.version = -1
	s.data = nil

	return false
}

func (s *EventStream) Error() error {
	return s.err
}

func (s *EventStream) Scan(event cqrs.Event) error {
	if s.version == -1 || s.data == nil {
		return ErrNoEventData
	}

	return json.Unmarshal(s.data, &event)
}
