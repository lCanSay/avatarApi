package main

import (
	//"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	//"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	models "github.com/lCanSay/avatarApi/pkg/models"
	"github.com/peterbourgon/ff/v3"

	//"github.com/gorilla/mux"
	//"github.com/joho/godotenv"

	//handler "github.com/lCanSay/avatarApi/api"
	//database "github.com/lCanSay/avatarApi/pkg/database"
	"github.com/lCanSay/avatarApi/pkg/jsonlog"

	//"github.com/lCanSay/avatarApi/pkg/models"
	_ "github.com/lib/pq"
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}

type application struct {
	config config
	models models.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	// This is a protected route
	// Only accessible to authenticated users with proper permissions
	fmt.Fprintf(w, "This is a protected route")
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	err := godotenv.Load("C:/KBTU/projectGo/avatarApi/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8080, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://testuser:1234@localhost:5432/avatar?sslmode=disable", "PostgreSQL DSN")
	)

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: models.NewModels(db),
		logger: logger,
	}

	if cfg.fill {
		// err = filler.PopulateDatabase(app.models)
		// if err != nil {
		// 	logger.PrintFatal(err, nil)
		// 	return
		// }
	}

	// Call app.server() to start the server.
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// https://github.com/golang-migrate/migrate?tab=readme-ov-file#use-in-your-go-project
	if cfg.migrations != "" {
		// driver, err := postgres.WithInstance(db, &postgres.Config{})
		// if err != nil {
		// 	return nil, err
		// }
		// m, err := migrate.NewWithDatabaseInstance(
		// 	cfg.migrations,
		// 	"postgres", driver)
		// if err != nil {
		// 	return nil, err
		// }
		// m.Up()
	}

	return db, nil
}
