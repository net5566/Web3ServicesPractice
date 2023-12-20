CREATE TABLE IF NOT EXISTS blocks (
    block_num INT PRIMARY KEY,
    block_hash VARCHAR(66) NOT NULL,
    block_time BIGINT NOT NULL,
    parent_hash VARCHAR(66) NOT NULL
);
