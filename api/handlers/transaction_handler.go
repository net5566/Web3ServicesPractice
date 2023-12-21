package handlers

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"web3-services/practice/constants"
	"web3-services/practice/types"
)

func normalizeHexString(txHash string) (string, bool) {
	if len(txHash) == 64 {
		return "0x" + txHash, true
	} else if len(txHash) == 66 && txHash[0] == '0' && txHash[1] == 'x' {
		return txHash, true
	}

	return "", false
}

func isValidHexString(txHash string) bool {
	match, _ := regexp.MatchString("^0x[0-9a-fA-F]+$", txHash)
	return match
}

func GetTransactionByHash(ginContext *gin.Context) {
	txHash := ginContext.Param("txHash")
	txHash, result := normalizeHexString(txHash)

	if result && isValidHexString(txHash) {
		mongodb, _ := ginContext.MustGet(constants.MONGO_DB).(*mongo.Database)
		transactionCollection := mongodb.Collection(constants.TRANSACTION_COLLECTION)
		filter := bson.M{"tx_hash": txHash}
		var result types.Transaction
		err := transactionCollection.FindOne(context.Background(), filter).Decode(&result)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				ginContext.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
				return
			}

			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Server Internal Error"})
		}

		ginContext.JSON(http.StatusOK, result)
	} else {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hexadecimal string"})
		return
	}
}
