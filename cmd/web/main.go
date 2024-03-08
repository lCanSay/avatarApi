package main

import (
	"context"
	"log"
	"net/http"

	// models "github.com/lCanSay/avatarApi/pkg/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	handler "github.com/lCanSay/avatarApi/api"
	database "github.com/lCanSay/avatarApi/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("C:/KBTU/projectGo/avatarApi/.env")
	if err != nil {
		log.Fatal("No .env file")
	}

	db := database.InitDB()
	defer db.Close()

	// later on will move this part
	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	router.HandleFunc("/health-check", handler.HealthCheck).Methods("GET")
	router.HandleFunc("/characters", handler.GetCharacters).Methods("GET")
	router.HandleFunc("/characters", handler.PostCharacter).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
