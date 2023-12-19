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

func main() {
	// load enviroment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	mongodbRootUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongodbRootPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongodbHost := os.Getenv("MONGO_HOST")
	mongodbPort := os.Getenv("MONGO_PORT")

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")

	// Connect to MySQL
	mysqlConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	mysqldb, err := gorm.Open(mysql.Open(mysqlConnectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	dbInstance, _ := mysqldb.DB()
	err = dbInstance.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")

	mongodbConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongodbRootUsername, mongodbRootPassword, mongodbHost, mongodbPort)
	mongodbClientOptions := options.Client().ApplyURI(mongodbConnectionString)

	// Connect to MongoDB
	mongodb, err := mongo.Connect(context.Background(), mongodbClientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongodb.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	defer func() {
		dbInstance, _ := mysqldb.DB()

		if err = dbInstance.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Disconnected from MySQL!")
	}()

	defer func() {
		if err = mongodb.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Disconnected from MongoDB!")
	}()
}
