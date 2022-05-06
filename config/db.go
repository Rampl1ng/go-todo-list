package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
	_ "github.com/joho/godotenv/autoload" // this help you to load .env automatically
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	logger, _     = thoth.Init("log")
	jsonLogger, _ = thoth.Init("json")
)

// MySQL
func Database() *sql.DB {
	user, exist := os.LookupEnv("DB_USER")
	if !exist {
		logger.Log(errors.New("DB_USER not set in .env"))
		log.Fatal("DB_USER not set in .env")
	}
	pass, exist := os.LookupEnv("DB_PASSWORD")
	if !exist {
		logger.Log(errors.New("DB_PASSWORD not set in .env"))
		log.Fatal("DB_PASSWORD not set in .env")
	}
	host, exist := os.LookupEnv("DB_HOST")
	if !exist {
		logger.Log(errors.New("DB_HOST not set in .env"))
		log.Fatal("DB_HOST not set in .env")
	}

	source := fmt.Sprintf("%s:%s@(%s:3306)/?charset=utf8&parseTime=True", user, pass, host)
	db, err := sql.Open("mysql", source)
	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}
	fmt.Println("Database Connection Successful")

	_, err = db.Exec(`CREATE DATABASE gotodo`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`USE gotodo`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE todos (
			id INT AUTO_INCREMENT,
			item TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			PRIMARY KEY (id)
		);
	`)
	if err != nil {
		fmt.Println(err)
	}

	return db
}

// MongoDB use json logging format
func MongoClient() (client *mongo.Client) {

	uri, exist := os.LookupEnv("MONGO_URI")
	if !exist {
		jsonLogger.Log(errors.New("MONGO_URI not set in .env"))
		log.Fatal("MONGO_URI not set in .env")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		jsonLogger.Log(err)
		log.Fatal(err)
	}
	fmt.Println("MongoDB Connection Successful")

	return client
}

func GetCollection() *mongo.Collection {
	client := MongoClient()

	db, exist := os.LookupEnv("MONGO_DB")
	if !exist {
		jsonLogger.Log(errors.New("MONGO_DB not set in .env"))
		log.Fatal("MONGO_DB not set in .env")
	}
	collection, exist := os.LookupEnv("MONGO_COLLECTION")
	if !exist {
		jsonLogger.Log(errors.New("MONGO_COLLECTION not set in .env"))
		log.Fatal("MONGO_COLLECTION not set in .env")
	}

	return client.Database(db).Collection(collection)
}
