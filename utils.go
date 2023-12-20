package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
}

type Block struct {
	BlockNum   int64  `gorm:"primaryKey;autoIncrement:false"`
	BlockHash  string `gorm:"not null"`
	BlockTime  int64  `gorm:"not null"`
	ParentHash string `gorm:"not null"`
}

func connectMySQL() *gorm.DB {
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")

	mysqlConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	mysqldb, err := gorm.Open(mysql.Open(mysqlConnectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")

	err = mysqldb.AutoMigrate(&Block{})

	if err != nil {
		log.Fatal(err)
	}

	return mysqldb
}

func connectMongoDB(mongoClient *mongo.Client) *mongo.Database {
	mongodbDatabase := os.Getenv("MONGO_INITDB_DATABASE")
	mongodb := mongoClient.Database(mongodbDatabase)

	return mongodb
}

func connectMongoClient() *mongo.Client {
	mongodbRootUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongodbRootPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongodbHost := os.Getenv("MONGO_HOST")
	mongodbPort := os.Getenv("MONGO_PORT")

	mongodbConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongodbRootUsername, mongodbRootPassword, mongodbHost, mongodbPort)
	clientOptions := options.Client().ApplyURI(mongodbConnectionString)

	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return mongoClient
}
