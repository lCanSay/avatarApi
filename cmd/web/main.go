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
	"github.com/joho/godotenv"
	handler "github.com/lCanSay/avatarApi/api"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("C:/KBTU/projectGo/avatarApi/.env")
	if err != nil {
		log.Fatal("No .env file")
	}

	db := InitDB()
	defer db.Close()

	// later will move this part
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

	log.Printf("Connect successful 1\n")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connect successful 2\n")
	//defer db.Close()

	return db
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
