CREATE TABLE IF NOT EXISTS Block (
    block_num INT PRIMARY KEY,
    block_hash VARCHAR(66) NOT NULL,
    block_time BIGINT NOT NULL,
    parent_hash VARCHAR(66) NOT NULL
);
