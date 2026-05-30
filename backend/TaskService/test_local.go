package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=taskuser password=taskpass123 dbname=tasksdb sslmode=disable"

	fmt.Println("Testing connection to local PostgreSQL...")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Open error:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping error:", err)
	}

	fmt.Println("✓ Successfully connected to local PostgreSQL!")

	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatal("Query error:", err)
	}

	fmt.Println("✓ Query result:", result)
}
