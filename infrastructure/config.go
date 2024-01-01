package infrastructure

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppEnv                 string
	ServerAddress          string
	ContextTimeout         int
	DBHost                 string
	DBPort                 string
	DBName                 string
	DBUser                 string
	DBPass                 string
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	timeout := 2

	contextTimeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	if err == nil {
		timeout = contextTimeout
	}

	tokenExpiry := 2
	refreshTokenExpiry := 16

	tx, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"))
	if err == nil {
		tokenExpiry = tx
	}

	rtx, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"))
	if err == nil {
		refreshTokenExpiry = rtx
	}

	return &Config{
		AppEnv:                 os.Getenv("APP_ENV"),
		ServerAddress:          os.Getenv("SERVER_ADDRESS"),
		ContextTimeout:         timeout,
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBName:                 os.Getenv("DB_NAME"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPass:                 os.Getenv("DB_PASS"),
		AccessTokenExpiryHour:  tokenExpiry,
		RefreshTokenExpiryHour: refreshTokenExpiry,
		AccessTokenSecret:      os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret:     os.Getenv("REFRESH_TOKEN_SECRET"),
	}
}
