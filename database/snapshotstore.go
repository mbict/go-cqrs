package database

import (
	"database/sql"
	"encoding/json"
	"github.com/mbict/go-cqrs/v4"
	"log"
	"time"
)

type snapshotStore struct {
	db         *sql.DB
	deleteStmt *sql.Stmt
	insertStmt *sql.Stmt
	selectStmt *sql.Stmt
}

func (s *snapshotStore) Load(aggregateId cqrs.AggregateId, aggregate cqrs.Aggregate) (int, error) {
	row := s.selectStmt.QueryRow(aggregateId, aggregate.AggregateName())

	var jsonData []byte
	var snapshotVersion int
	err := row.Scan(&jsonData, &snapshotVersion)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return snapshotVersion, json.Unmarshal(jsonData, aggregate)
}

func (s *snapshotStore) Write(aggregate cqrs.Aggregate) error {
	payload, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	_, err = tx.Stmt(s.deleteStmt).Exec(aggregate.AggregateId(), aggregate.AggregateName(), aggregate.Version())
	if err != nil {
		return err
	}

	result, err := tx.Stmt(s.insertStmt).Exec(aggregate.AggregateId(), aggregate.AggregateName(), payload, aggregate.Version(), time.Now())
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func NewPostgresSnapshotStore(db *sql.DB) cqrs.SnapshotStore {

	deleteStmt, err := db.Prepare("DELETE FROM snapshots WHERE aggregate_id = $1 AND aggregate_type = $2 AND version <> $3")
	if err != nil {
		log.Fatal(err)
	}

	insertStmt, err := db.Prepare("INSERT INTO snapshots (aggregate_id, aggregate_type, data, version, created) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}

	selectStmt, err := db.Prepare("SELECT data, version FROM snapshots WHERE aggregate_id = $1 AND aggregate_type = $2 ORDER BY version DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}

	return &snapshotStore{
		db:         db,
		deleteStmt: deleteStmt,
		insertStmt: insertStmt,
		selectStmt: selectStmt,
	}
}

func NewMySQLSnapshotStore(db *sql.DB) cqrs.SnapshotStore {

	deleteStmt, err := db.Prepare("DELETE FROM snapshots WHERE aggregate_id = ? AND aggregate_type = ? AND version <> ?")
	if err != nil {
		log.Fatal(err)
	}

	insertStmt, err := db.Prepare("INSERT INTO snapshots (aggregate_id, aggregate_type, data, version, created) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	selectStmt, err := db.Prepare("SELECT data, version FROM snapshots WHERE aggregate_id = ? AND aggregate_type = ? ORDER BY version DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}

	return &snapshotStore{
		db:         db,
		deleteStmt: deleteStmt,
		insertStmt: insertStmt,
		selectStmt: selectStmt,
	}
}
