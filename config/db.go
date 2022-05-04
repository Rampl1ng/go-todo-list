package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
	_ "github.com/joho/godotenv/autoload" // t help you to load .env automatically
)

func Database() *sql.DB {
	logger, _ := thoth.Init("log")

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
	} else {
		fmt.Println("Database Connection Successful")
	}

	_, err = db.Exec(`CREATE DATABASE gotodo`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`USE gotodo`)
	if err != nil {
		fmt.Println(err)
	}

	db.Exec(`
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
