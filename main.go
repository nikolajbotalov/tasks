package main

import (
	router "TaskFlow/internal/delivery/http"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
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

	r := router.AppRouters(db)

	r.Run()
}
