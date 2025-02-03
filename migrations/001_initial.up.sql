-- Postgres
CREATE TABLE IF NOT EXISTS rates (
    timestamp BIGINT NOT NULL,
    ask DOUBLE PRECISION NOT NULL,
    bid DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS rates_created_at_idx ON rates (created_at);



-- ClickHouse
-- CREATE TABLE IF NOT EXISTS rates (
--     timestamp UInt64,
--     ask Float64,
--     bid Float64,
--     created_at DateTime DEFAULT now()
-- ) ENGINE = MergeTree()
-- ORDER BY created_at;