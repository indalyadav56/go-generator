package main

import (
	"backend/config"
	"fmt"
	"log"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Printf("Error creating application: %v\n", err)
		return
	}

	if err := application.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		return
	}

	defer func() {
		if err := application.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down application: %v\n", err)
		}
	}()
}
