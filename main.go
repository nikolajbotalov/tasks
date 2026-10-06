package main

import (
	"TaskFlow/internal/config"
	router "TaskFlow/internal/delivery/http"
	authRepo "TaskFlow/internal/repository/auth"
	"TaskFlow/internal/repository/postgres"
	tasksRepo "TaskFlow/internal/repository/tasks"
	authUC "TaskFlow/internal/use-cases/auth"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=taskflow password=admin sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully connected to db")

	authRepository := authRepo.NewAuthRepository(db)
	authUseCases := authUC.NewAuthUseCases(authRepository, cfg)
	taskRepository := tasksRepo.NewTasksRepository(db)
	tasksUseCases := tasksUC.NewTasksUseCases(taskRepository)

	r := router.AppRouters(authUseCases, cfg, tasksUseCases)

	ctx := context.Background()

	if err = postgres.EnsureSchema(ctx, db, "tasks", "users"); err != nil {
		log.Printf("schema check failed: %v\n", err)
	}

	r.Run()
}
