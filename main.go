package main

import (
	"errors"
	"log"

	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"github.com/rampl1ng/go-todoList/routes"
)

func main() {
	logger, _ := thoth.Init("log")

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}
	log.Println(".env file load successfully")

	// port, exist := os.LookupEnv("PORT")
	// if !exist {
	// 	logger.Log(errors.New("PORT not set in .env"))
	// 	log.Fatal("PORT not set in .env")
	// }

	r := routes.SetUpRouters()
	err := r.Run()
	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}

	// err := http.ListenAndServe(":"+port, routes.Init())
	// if err != nil {
	// 	logger.Log(err)
	// 	log.Fatal(err)
	// }
}
