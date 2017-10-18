CREATE TABLE events
(
    id INTEGER PRIMARY KEY NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    data JSON,
    version INTEGER NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX events_aggregate_id_aggregate_type_version_index ON events (aggregate_id, aggregate_type, version);
CREATE SEQUENCE events_id_seq;
ALTER TABLE events ALTER id SET DEFAULT NEXTVAL('events_id_seq');
SELECT setval('events_id_seq', 1, false);