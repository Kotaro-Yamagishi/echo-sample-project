package sqlboiler

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type SQLBoilerImpl struct {
}

func NewSQLBoilerImpl() SQLBoiler {
	sqlboiler := new(SQLBoilerImpl)
	return sqlboiler
}

func (s *SQLBoilerImpl) ConnectDB() *sql.DB {
	// Get database connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "sakila")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "passw0rd")

	// Create connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
