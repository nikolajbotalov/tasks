package main

import (
	"TaskFlow/internal/app"

	_ "TaskFlow/docs"

	_ "github.com/lib/pq"
)

// @title Swagger TaskFlow API
// @version 1.0
// @description This is a simple task flow server

// @BasePath /api/v1

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer application.Close()

	application.Server.Run()
}
