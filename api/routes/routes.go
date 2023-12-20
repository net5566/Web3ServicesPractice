package routes

import (
	"github.com/gin-gonic/gin"

	"web3-services/practice/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	blocks := r.Group("/blocks")
	{
		blocks.GET("/", handlers.GetLatestBlocks)
		blocks.GET("/:blockNumber", handlers.GetBlockByNumber)
	}

	r.GET("/transaction/:txHash", handlers.GetTransactionByHash)
}
