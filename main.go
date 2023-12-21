package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web3-services/practice/api/routes"
	"web3-services/practice/services"

	"github.com/gin-gonic/gin"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	loadDotEnv()
	mysqldb := connectMySQL()
	mongoClient := connectMongoClient()
	mongodb := connectMongoDB(mongoClient)
	rpcClient := establishRPC()

	go func() {
		// Run per 15 seconds
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			services.IndexBlockRPC(mongodb, mysqldb, rpcClient, 0)
		}
	}()

	router := gin.Default()
	routes.SetupRoutes(router)
	go router.Run("localhost:8080")

	select {
	case <-sigChan:
		fmt.Println("Received interrupt signal. Cleaning up...")
		handleMySQLDisconnected(mysqldb)
		handleMongoDisconnected(mongoClient)
		os.Exit(0)
	}
}
