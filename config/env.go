package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string `env:"PUBLIC_HOST"`
	Port                   string `env:"PORT"`
	DBUser                 string `env:"DB_USER"`
	DBPassword             string `env:"DB_PASSWORD"`
	DBAddress              string `env:"DB_ADDRESS"`
	DBName                 string `env:"DB_NAME"`
	JWTSecret              string `env:"JWT_SECRET"`
	JWTExpirationInSeconds int64  `env:"JWT_EXPIRATION"`
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "admin"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "products"),
		JWTSecret:              getEnv("JWT_SECRET", "custom_secret_key"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7), // 7 days
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
