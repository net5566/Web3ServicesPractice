package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"web3-services/practice/api/routes"
	"web3-services/practice/services"
)

func main() {
	loadDotEnv()
	mysqldb := connectMySQL()
	mongoClient := connectMongoClient()
	mongodb := connectMongoDB(mongoClient)
	rpcClient := establishRPC()

	defer handleMySQLDisconnected(mysqldb)
	defer handleMongoDisconnected(mongoClient)

	go func() {
		// Run per 15 seconds
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			services.IndexBlockRPC(mongodb, mysqldb, rpcClient, 0)
		}
	}()

	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run("localhost:8080")

	select {}
}
