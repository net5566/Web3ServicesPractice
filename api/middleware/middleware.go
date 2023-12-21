package middleware

import (
	"web3-services/practice/constants"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func AttachDatabasesToContext(mongodb *mongo.Database, mysqldb *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(constants.MONGO_DB, mongodb)
		context.Set(constants.MYSQL_DB, mysqldb)
		context.Next()
	}
}
