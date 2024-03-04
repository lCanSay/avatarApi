package main

import (
	//"encoding/json"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// models "github.com/lCanSay/avatarApi/pkg/models"

	"github.com/gorilla/mux"
	handler "github.com/lCanSay/avatarApi/api"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/characters", handler.GetCharacters).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func InitDB() *sql.DB {
	user := getEnv("DB_USER")
	password := getEnv("DB_PASSWORD")
	dbname := getEnv("DB_NAME")

	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s port=5432", user, dbname, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return db
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
