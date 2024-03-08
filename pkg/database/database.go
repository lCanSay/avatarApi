package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	user := getEnv("DB_USER")
	password := getEnv("DB_PASSWORD")
	dbname := getEnv("DB_NAME")

	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s port=5432 sslmode=disable", user, dbname, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	log.Printf("Connect successful 1\n")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connect successful 2\n")

	return db
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
