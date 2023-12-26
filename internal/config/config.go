package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GetConfig()
}

var (
	c *Config
)

type Config struct {
	PostgreSQL PostgreSQL
}

type PostgreSQL struct {
	Port     string
	User     string
	Password string
	DBName   string
}

func GetConfig() *Config {
	if c == nil {
		//REDIS

		port := os.Getenv("DB_PORT")
		if port == "" {
			panic("PORT is not set")
		}

		user := os.Getenv("DB_USERNAME")
		if user == "" {
			panic("USER is not set")
		}

		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			panic("PASSWORD is not set")
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			panic("DBNAME is not set")
		}

		c = &Config{
			PostgreSQL: PostgreSQL{
				Port:     port,
				User:     user,
				Password: password,
				DBName:   dbName,
			},
		}

		return c
	}

	return c
}
