# Web3 Services Practice

## Run Databases
```bash
$ docker compose up -d
```

## Get Go Dependencies
```bash
$ go get .
```

## Run Go Server
```bash
$ go run .
```

## Schema
Please check `init.sql` and `inti-mongo.js`

Splitting data into three Tables

### Block Data

The Block Data are put in the MySQL Database.

It provides ACID properties and various query options.

The RPC provides many kinds of queries for Blocks. 

Therefore, it is better to store Block Data into Relational Databases for further development.

### Mapping from Block Number to Corresponding TxHashes

And then, the mappings from Block Number to Transaction Hashses are placed in MongoDB.

MongoDB is suitable for storing array data.

Splitting the data, because `GetLatestBlocks` does not need to respond with the mapping.

By doing so, we can efficiently reduce the size of single table.

Also, the speed of queries of MongoDB is much faster than MySQL.

`GetBlockByBlockNumber` will not be impacted if it needs one more MongoDB query.

### Transaction Data

Last, the Transactions Data are stored on MongoDB.

As previously mentioned, the `GetTransactionByHash` can be extremely fast with only NoSQL query.

Following the documents of RPC-JSON API, the querying APIs for Transaction Data are more simple.

Looking up Transaction Data with TxHash is compatible for the APIs querying for Transaction Data.

## Ethereum Block Indexer

|               | Block Data | BlockTransactions | Transaction Data |
| :------------ | :--------: | :---------------: | :--------------: |
| Database      |   MySQL    |       Mongo       |      Mongo       |
| Batch Size    |     20     |      Various      |        1         |
| Upsert        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |

### Basic Config

Every 15 seconds, `Ethereum Block Indexer` will be executed once.

For each indexing operation, at most 100 blocks will be updated/inserted.

### Upsert Feature

`Upsert` feature is activated in all the inserting operations for the robustness of indexing.

### Block Stability Update

Due to the stability issue of blocks, the indexing starts from the last 20th stored block.
