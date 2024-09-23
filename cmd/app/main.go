package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/debate-io/service-auth/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	config, err := app.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(config)
	app.Initialize()
	app.RunApp()
}
