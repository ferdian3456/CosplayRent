package app

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"time"
)

func NewDB() *sql.DB {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath := filepath.Join(cwd, "../.env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal(err)
	}

	dbUri := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("pgx", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping to database:%v\n", err)
		log.Fatal(err)
	}

	log.Println("Success to connect and ping database")
	log.Println("Server is running on port 8081")

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
