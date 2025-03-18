package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Global variable to hold the pool connection
var Pool *pgxpool.Pool

// InitDB function to initialize the PostgreSQL connection pool
func InitDB(username, password, host, port, dbname string) error {
	// สร้าง connection string ที่รวม username, password, host, dbname, และ port
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)

	// สร้าง connection pool โดยใช้ pgxpool.New
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Database connection pool initialized")
	return nil
}

// CloseDB function to close the connection pool
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		fmt.Println("Database connection pool closed")
	}
}

// PingDB function to test the connection
func PingDB() error {
	// Use a context with timeout to ping the DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Pool.Ping(ctx)
}
