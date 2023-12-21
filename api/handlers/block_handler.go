package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"web3-services/practice/constants"
	"web3-services/practice/types"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func GetBlockByNumber(ginContext *gin.Context) {
	blockNumStr := ginContext.Param("blockNumber")
	blockNum, err := strconv.Atoi(blockNumStr)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blcok number"})
	}

	var block types.Block
	mysqldb, _ := ginContext.MustGet(constants.MYSQL_DB).(*gorm.DB)
	err = mysqldb.Where("block_num = ?", blockNum).First(&block).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
			return
		}
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Server Internal Error"})
		return
	}

	mongodb, _ := ginContext.MustGet(constants.MONGO_DB).(*mongo.Database)
	blockTransactionsCollection := mongodb.Collection(constants.BLOCK_TRANSACTINOS_COLLECTION)
	filter := bson.M{"block_num": blockNum}
	var blockTransactions types.BlockTransactions
	err = blockTransactionsCollection.FindOne(context.Background(), filter).Decode(&blockTransactions)

	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Server Internal Error"})
		}
	}

	result := types.BlockComplete{
		BlockNum:     block.BlockNum,
		BlockHash:    block.BlockHash,
		BlockTime:    block.BlockTime,
		ParentHash:   block.ParentHash,
		Transactions: blockTransactions.Transactions,
	}

	ginContext.JSON(http.StatusOK, result)
}

func GetLatestBlocks(ginContext *gin.Context) {
	limitStr := ginContext.Query("limit")
	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
	}

	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}

	mysqldb, _ := ginContext.MustGet(constants.MYSQL_DB).(*gorm.DB)

	var blocks []types.Block
	err = mysqldb.Order("block_num DESC").Limit(limit).Find(&blocks).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
			return
		}
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Server Internal Error"})
		return
	}

	ginContext.JSON(http.StatusOK, blocks)
}
