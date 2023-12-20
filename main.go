package main

import (
	"fmt"
)

func main() {
	loadDotEnv()
	mysqldb := connectMySQL()
	mongoClient := connectMongoClient()

	rpcClient := establishRPC()
	blockHeight := getBlockHeight(rpcClient)
	fmt.Printf("Current block height: %d\n", blockHeight)

	defer handleMySQLDisconnected(mysqldb)
	defer handleMongoDisconnected(mongoClient)
}
