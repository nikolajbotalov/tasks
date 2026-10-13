package main

import (
	"TaskFlow/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer application.Close()

	application.Server.Run()
}
