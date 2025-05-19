package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	dbHostEnv     = "DB_HOST"
	dbPortEnv     = "DB_PORT"
	dbNameEnv     = "DB_NAME"
	dbUserEnv     = "DB_USER"
	dbPasswordEnv = "DB_PASSWORD"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfig() *DBConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	return &DBConfig{
		DBHost:     getEnv(dbHostEnv, "Db host error"),
		DBPort:     getEnv(dbPortEnv, "Db port error"),
		DBName:     getEnv(dbNameEnv, "Db name error"),
		DBUser:     getEnv(dbUserEnv, "Db user error"),
		DBPassword: getEnv(dbPasswordEnv, "Db password error"),
	}
}

func getEnv(key, errMsg string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		log.Fatal(errMsg)
	}

	return value
}
