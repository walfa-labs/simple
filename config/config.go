package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Mongo    MongoConfig
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
	Enabled  bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Enabled  bool
}

type MongoConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	DB         string
	AuthSource string
	Enabled    bool
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", ""),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DB:       getEnv("POSTGRES_DB", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Mongo: MongoConfig{
			Host:       getEnv("MONGO_HOST", ""),
			Port:       getEnvInt("MONGO_PORT", 27017),
			User:       getEnv("MONGO_USER", ""),
			Password:   getEnv("MONGO_PASSWORD", ""),
			DB:         getEnv("MONGO_DB", ""),
			AuthSource: getEnv("MONGO_AUTH_SOURCE", ""),
		},
	}

	// Enable integrations only if required env vars are set
	cfg.Postgres.Enabled = cfg.Postgres.Host != "" && cfg.Postgres.User != "" && cfg.Postgres.DB != ""
	cfg.Redis.Enabled = cfg.Redis.Host != ""
	cfg.Mongo.Enabled = cfg.Mongo.Host != "" && cfg.Mongo.DB != ""

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DB,
	)
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *MongoConfig) ConnectionString() string {
	authSource := c.AuthSource
	if authSource == "" {
		authSource = "admin" // default to admin for root user auth
	}
	if c.User != "" && c.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			c.User, c.Password, c.Host, c.Port, c.DB, authSource)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s",
		c.Host, c.Port, c.DB)
}
