CREATE TABLE snapshots
(
    id SERIAL PRIMARY KEY NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    data JSON,
    version INTEGER NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX snapshots_aggregate_id_aggregate_type_version_index ON snapshots (aggregate_id, aggregate_type, version);