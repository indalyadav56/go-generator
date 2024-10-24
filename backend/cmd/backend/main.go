package main

import (
	"backend/config"
	"fmt"
	"log"
)

func main() {
	app, err := config.InitializeApp(".env")
	if err != nil {
		panic(err)
	}

	appRouter := config.SetupRouter(app)

	if err := appRouter.Run(fmt.Sprintf(":%s", app.Config.ServerPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
