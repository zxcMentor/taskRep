package main

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/handler"
	"awesomeProject/internal/repository"
	"awesomeProject/router"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// @title API Title
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	dbHost := "db"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "password"
	dbName := "postgres"
	sslmode := "disable"

	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + sslmode

	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		log.Fatalf("Could not ping to the database: %v", err)
	}

	userRepository := repository.NewDBPostgres(database)

	userHandler := handler.NewUserHandler(userRepository)

	r := router.SetupRouter(userHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
