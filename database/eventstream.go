package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/mbict/go-cqrs/v4"
	"time"
)

var (
	ErrNoEventData = errors.New("Cannot scan, no event data")
)

type EventStream struct {
	rows *sql.Rows

	err         error
	version     int
	aggregateId cqrs.AggregateId
	eventType   cqrs.EventType
	data        sql.RawBytes
	timestamp   time.Time
}

func NewDatabaseEventStream(rows *sql.Rows) cqrs.EventStream {
	return &EventStream{
		rows: rows,
	}
}

func (s *EventStream) AggregateId() cqrs.AggregateId {
	return s.aggregateId
}

func (s *EventStream) EventType() cqrs.EventType {
	return s.eventType
}

func (s *EventStream) Version() int {
	return s.version
}

func (s *EventStream) Timestamp() time.Time {
	return s.timestamp
}

func (s *EventStream) Next() bool {
	if s.rows.Next() {
		if s.err = s.rows.Scan(s.aggregateId, &s.eventType, &s.data, &s.version, &s.timestamp); s.err != nil {
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

func (s *EventStream) Scan(event cqrs.EventData) error {
	if s.version == -1 || s.data == nil {
		return ErrNoEventData
	}

	return json.Unmarshal(s.data, &event)
}
