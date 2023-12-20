package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	loadDotEnv()
	mysqldb := connectMySQL()
	mongoClient := connectMongoClient()

	rpcClient := establishRPC()
	blockHeight := getBlockHeight(rpcClient)
	fmt.Printf("Current block height: %d\n", blockHeight)

	defer func() {
		dbInstance, _ := mysqldb.DB()

		if err := dbInstance.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Disconnected from MySQL!")
	}()

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Disconnected from MongoDB!")
	}()
}
