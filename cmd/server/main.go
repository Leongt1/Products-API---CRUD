package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"main.go/internal/handler"
)

func main() {
	connStr := "host=localhost user=noel password=noelgt28 dbname=productsdb sslmode=disable" // Corrected format
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // It's good practice to close the database connection when it's no longer needed.

	router := gin.Default()

	productHandler := handler.NewProductHandler(db)
	productHandler.RegisterProductRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run the server: ", err)
	}
}
