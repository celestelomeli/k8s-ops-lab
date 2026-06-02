// db.go manages the database connection pool
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib" // registers pgx as the postgres driver
)

// DB is the shared connection used by all handlers
var DB *sql.DB

// initDB opens the connection pool and verifies connectivity
func initDB() {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    // retry up to 10 times, 2 seconds between attempts
    for i := 0; i < 10; i++ {
        DB, err = sql.Open("pgx", dsn)
        if err == nil {
            err = DB.Ping() // Ping tests the connection
        }
        if err == nil {
            log.Println("connected to database")
            return
        }
        log.Printf("database not ready, retrying in 2s (%d/10)...", i+1)
        time.Sleep(2 * time.Second)
    }
    log.Fatal("could not connect to database: ", err)
}