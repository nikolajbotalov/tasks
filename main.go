package main

import (
	"TaskFlow/internal/config"
	router "TaskFlow/internal/delivery/http"
	authRepository "TaskFlow/internal/repository/auth"
	"TaskFlow/internal/repository/postgres"
	authUC "TaskFlow/internal/use-cases/auth"
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

	auth := authRepository.NewAuthRepository(db)
	authUseCases := authUC.NewAuthUseCases(auth, cfg)

	r := router.AppRouters(db, authUseCases)

	ctx := context.Background()

	if err = postgres.EnsureSchema(ctx, db, "tasks", "users"); err != nil {
		log.Printf("schema check failed: %v\n", err)
	}

	r.Run()
}
