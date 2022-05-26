package main

import (
	"errors"
	"log"

	"github.com/joho/godotenv"
	"github.com/rampl1ng/go-todoList/routes"
	"github.com/rampl1ng/go-todoList/utils"
)

func main() {
	logger, _ := utils.Init("log")

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}
	log.Println(".env file load successfully")

	r := routes.SetUpRouters()

	err := r.Run()
	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}
}
