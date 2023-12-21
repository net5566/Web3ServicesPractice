package routes

import (
	"github.com/gin-gonic/gin"

	"web3-services/practice/api/handlers"
)

func SetupRoutes(router *gin.Engine) {
	blocks := router.Group("/blocks")
	{
		blocks.GET("/", handlers.GetLatestBlocks)
		blocks.GET("/:blockNumber", handlers.GetBlockByNumber)
	}

	router.GET("/transaction/:txHash", handlers.GetTransactionByHash)
}
