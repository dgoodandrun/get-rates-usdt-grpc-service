CREATE TABLE IF NOT EXISTS rates (
    timestamp UInt64,
    ask Float64,
    bid Float64,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY created_at;