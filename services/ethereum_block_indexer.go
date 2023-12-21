package services

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"web3-services/practice/types"
)

func convertToBlock(block *RPCBlock) types.Block {
	blockNum, _ := strconv.ParseInt(block.Number, 0, 64)
	blockTime, _ := strconv.ParseInt(block.Timestamp, 0, 64)

	return types.Block{
		BlockNum:   blockNum,
		BlockHash:  block.Hash,
		BlockTime:  blockTime,
		ParentHash: block.ParentHash,
	}
}

// Get the height of blocks from RPC
func getBlockHeight(rpcClient *rpc.Client) int64 {
	var blockHeightString string
	err := rpcClient.Call(&blockHeightString, "eth_blockNumber")

	if err != nil {
		log.Fatal(err)
	}

	blockHeight, err := strconv.ParseInt(blockHeightString, 0, 64)

	if err != nil {
		log.Fatal(err)
	}

	return blockHeight
}

func getBlockByNumber(rpcClient *rpc.Client, blockNum int) *RPCBlock {
	var block RPCBlock
	err := rpcClient.Call(&block, "eth_getBlockByNumber", fmt.Sprintf("0x%x", blockNum), true)

	if err != nil {
		log.Fatal(err)
	}

	return &block
}

func getTransactionReceipt(rpcClient *rpc.Client, txHash *string) *RPCTransactionReceipt {
	var transactionReceipt RPCTransactionReceipt

	err := rpcClient.Call(&transactionReceipt, "eth_getTransactionReceipt", *txHash)

	if err != nil {
		log.Fatal(err)
	}

	return &transactionReceipt
}

func insertBlockBatch(mysqldb *gorm.DB, blockBatch []types.Block, batchSize int) {
	result := mysqldb.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(blockBatch, batchSize)

	if result.Error != nil {
		panic(result.Error)
	}
}

func insertBlockTransactions(blockTransactionsCollection *mongo.Collection, blockTransactions types.BlockTransactions) {
	// Prepare upsert operations
	upsertFilter := bson.D{{"block_num", blockTransactions.BlockNum}}
	upsertUpdate := bson.D{
		{"$set", bson.D{{"transactions", blockTransactions.Transactions}}},
		{"$setOnInsert", bson.D{{"block_num", blockTransactions.BlockNum}}},
	}

	// Execute update one with upsert operations
	_, err := blockTransactionsCollection.UpdateOne(context.Background(), upsertFilter, upsertUpdate, options.Update().SetUpsert(true))
	if err != nil {
		panic(err)
	}
}

func insertTransactions(transactionCollection *mongo.Collection, transactionInterfaces []interface{}) {
	// Prepare upsert operations
	var upsertOperations []mongo.WriteModel
	for _, transaction := range transactionInterfaces {
		upsertFilter := bson.D{{"tx_hash", transaction.(types.Transaction).Hash}}
		upsertUpdate := bson.D{
			{"$set", bson.D{
				{"from", transaction.(types.Transaction).From},
				{"to", transaction.(types.Transaction).To},
				{"value", transaction.(types.Transaction).Value},
				{"nonce", transaction.(types.Transaction).Nonce},
				{"data", transaction.(types.Transaction).Data},
				{"logs", transaction.(types.Transaction).Logs},
			}},
			{"$setOnInsert", bson.D{{"tx_hash", transaction.(types.Transaction).Hash}}},
		}

		// Create upsert operation
		upsertOperation := mongo.NewUpdateOneModel()
		upsertOperation.SetFilter(upsertFilter)
		upsertOperation.SetUpdate(upsertUpdate)
		upsertOperation.SetUpsert(true)

		// Append upsert operation to the list
		upsertOperations = append(upsertOperations, upsertOperation)
	}

	// Execute bulk write with upsert operations
	_, err := transactionCollection.BulkWrite(context.Background(), upsertOperations, options.BulkWrite().SetOrdered(false))
	if err != nil {
		panic(err)
	}
}

func integrateTransactionData(transactionRPCData *RPCTransaction, transactionReceipt *RPCTransactionReceipt) types.Transaction {
	nonce, _ := strconv.ParseInt(transactionRPCData.Nonce, 0, 64)

	transaction := types.Transaction{
		Hash:  transactionRPCData.Hash,
		From:  transactionRPCData.From,
		To:    transactionRPCData.To,
		Value: transactionRPCData.Value,
		Nonce: int(nonce),
		Data:  transactionRPCData.Data,
		Logs:  []types.Log{},
	}

	if len(transactionReceipt.Logs) > 0 {
		for _, logData := range transactionReceipt.Logs {
			index, err := strconv.ParseInt(logData.Index, 0, 64)

			if err != nil {
				log.Fatal(err)
			}

			transaction.Logs = append(transaction.Logs, types.Log{
				Data:  logData.Data,
				Index: int(index),
			})
		}
	}

	return transaction
}

func initParams(mysqldb *gorm.DB, rpcClient *rpc.Client, initBlockNum int) (int, int) {
	var count int64
	var startNum int
	var endNum int

	if err := mysqldb.Table("blocks").Count(&count).Error; err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		var lastBlock types.Block
		mysqldb.Last(&lastBlock)
		startNum = int(lastBlock.BlockNum) + 1
		blockHeight := getBlockHeight(rpcClient)

		if blockHeight-lastBlock.BlockNum < 100 {
			endNum = int(blockHeight)
		} else {
			endNum = startNum + 99
		}
	} else {
		startNum = initBlockNum
		endNum = startNum + 99
	}

	return startNum, endNum
}

func IndexBlockRPC(mongodb *mongo.Database, mysqldb *gorm.DB, rpcClient *rpc.Client, initBlockNum int) {
	const batchSize = 20
	blockCount := 0
	var blockBatch []types.Block

	startNum, endNum := initParams(mysqldb, rpcClient, initBlockNum)

	for blockNum := startNum; blockNum <= endNum; blockNum++ {
		block := getBlockByNumber(rpcClient, blockNum)
		blockBatch = append(blockBatch, convertToBlock(block))
		blockCount += 1

		if blockCount == batchSize {
			// Avoid concurrent modification
			blockBatchCopy := make([]types.Block, batchSize)
			copy(blockBatchCopy, blockBatch)

			go insertBlockBatch(mysqldb, blockBatchCopy, batchSize)

			blockCount = 0
			blockBatch = nil
		}

		transactionCollection := mongodb.Collection("Transaction")
		blockTransactionsCollection := mongodb.Collection("BlockTransactions")

		if len(block.Transactions) > 0 {
			var transactionHashes []string
			var transactionInterfaces []interface{}

			for _, transactionRPCData := range block.Transactions {
				transactionReceipt := getTransactionReceipt(rpcClient, &transactionRPCData.Hash)
				transaction := integrateTransactionData(&transactionRPCData, transactionReceipt)

				transactionInterfaces = append(transactionInterfaces, transaction)
				transactionHashes = append(transactionHashes, transaction.Hash)
			}

			blockTransactions := types.BlockTransactions{
				BlockNum:     blockNum,
				Transactions: transactionHashes,
			}

			go insertTransactions(transactionCollection, transactionInterfaces)
			go insertBlockTransactions(blockTransactionsCollection, blockTransactions)
		}
	}

	if blockCount > 0 {
		insertBlockBatch(mysqldb, blockBatch, blockCount)
		blockBatch = nil
	}
}
