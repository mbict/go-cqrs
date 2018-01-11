CREATE TABLE snapshots
(
  id BIGINT(20) UNSIGNED PRIMARY KEY NOT NULL AUTO_INCREMENT,
  aggregate_id VARBINARY(16) NOT NULL,
  aggregate_type VARCHAR(255) NOT NULL,
  data BLOB,
  version INT(11) UNSIGNED NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  INDEX snapshots_aggregate_type_aggregate_id_index (aggregate_type, aggregate_id, version)
);
