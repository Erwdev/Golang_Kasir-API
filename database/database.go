package database

import (
	"context"
	"database/sql"
	"log"
	"time"
	_ "github.com/lib/pq"
)

func InitDB(connectionString string )(*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}


	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	log.Println("Database connected successfully")
	return db , nil
}

func HealthCheck(db *sql.DB) error{
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}