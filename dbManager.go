package main

import (
	"fmt"
	"log"
	// "os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id                 int
	First, Last, Email string
}

func NewDBManager() *sqlx.DB {
	dbUser := "user"
	// dbUser := os.Getenv("DB_USER")
	dbPassword := "password"
	// dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := "3306"
	// dbPort := os.Getenv("DB_PORT")
	dbName := "stock_price"
	// dbName := os.Getenv("DB_NAME")
	dbHost := "localhost"
	// dbHost := os.Getenv("DB_HOST")
	key := fmt.Sprintf("%s:%s@(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sqlx.Connect("mysql", key)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
