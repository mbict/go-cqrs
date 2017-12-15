package memory

import (
	"fmt"
	"github.com/mbict/go-cqrs"
	"github.com/satori/go.uuid"
	"github.com/square/go-jose/json"
	"sync"
)

type snapshot struct {
	version int
	data    []byte
}

type snapshotStore struct {
	snapshots map[string]snapshot
	mu        sync.RWMutex
}

func (s *snapshotStore) Load(aggregateId uuid.UUID, aggregate cqrs.Aggregate) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := aggregate.AggregateName() + "-" + aggregateId.String()
	snapshot, ok := s.snapshots[key]
	if !ok {
		return 0, nil
	}

	fmt.Println("load from snapshot")
	return snapshot.version, json.Unmarshal(snapshot.data, aggregate)
}

func (s *snapshotStore) Write(aggregate cqrs.AggregateComposition) error {
	data, err := json.Marshal(aggregate.(cqrs.Aggregate))
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := aggregate.AggregateName() + "-" + aggregate.AggregateId().String()
	s.snapshots[key] = snapshot{
		version: aggregate.Version(),
		data:    data,
	}

	fmt.Println("store snapshot")
	return nil
}

func NewSnapshotStore() cqrs.SnapshotStore {
	return &snapshotStore{
		snapshots: make(map[string]snapshot),
	}
}
