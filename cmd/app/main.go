package main

import (
	"golang-gin-app/internal/app"
)

func main() {
	appInstance := app.NewApp()

	// Start the server
	if err := appInstance.Run(""); err != nil {
		panic(err)
	}
}
