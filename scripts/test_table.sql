CREATE TABLE IF NOT EXISTS test.metrics ON CLUSTER '{cluster}'
(
    `timestamp`           DateTime64(9,'Asia/Shanghai') CODEC (DoubleDelta),
    `metric_group`        LowCardinality(String),
    `number_field_keys`   Array(LowCardinality(String)),
    `number_field_values` Array(Float64),
    `string_field_keys`   Array(LowCardinality(String)),
    `string_field_values` Array(String),
    `tag_keys`            Array(LowCardinality(String)),
    `tag_values`          Array(LowCardinality(String)),
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{cluster}-{shard}/{database}/metrics', '{replica}')
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (metric_group, timestamp)
TTL toDateTime(timestamp) + INTERVAL 7 DAY;

CREATE TABLE IF NOT EXISTS test.metrics_all ON CLUSTER '{cluster}' AS test.metrics
ENGINE = Distributed('{cluster}', test, metrics, rand());