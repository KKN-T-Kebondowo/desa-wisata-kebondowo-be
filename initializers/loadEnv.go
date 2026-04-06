package initializers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBUserName     string
	DBUserPassword string
	DBName         string
	DBPort         string
	ServerPort     string

	ClientOrigin string

	AccessTokenPrivateKey  string
	AccessTokenPublicKey   string
	RefreshTokenPrivateKey string
	RefreshTokenPublicKey  string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	AccessTokenMaxAge      int
	RefreshTokenMaxAge     int
}

func LoadConfig(path string) (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("loading .env file: %w", err)
	}

	accessTokenExpInHours, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRED_IN_HOURS"))
	if err != nil {
		accessTokenExpInHours = 1 // default to 1 hour
	}

	refreshTokenExpInDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRED_IN_DAYS"))
	if err != nil {
		refreshTokenExpInDays = 7 // default to 7 days
	}

	config := Config{
		DBHost:         os.Getenv("POSTGRES_HOST"),
		DBUserName:     os.Getenv("POSTGRES_USER"),
		DBUserPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		DBPort:         os.Getenv("POSTGRES_PORT"),
		ServerPort:     os.Getenv("PORT"),

		ClientOrigin: os.Getenv("CLIENT_ORIGIN"),

		AccessTokenPrivateKey: os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"),
		AccessTokenPublicKey:  os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"),

		RefreshTokenPrivateKey: os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"),
		RefreshTokenPublicKey:  os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"),

		AccessTokenExpiresIn:  time.Duration(accessTokenExpInHours) * time.Hour,
		RefreshTokenExpiresIn: time.Duration(refreshTokenExpInDays) * 24 * time.Hour,

		AccessTokenMaxAge:  accessTokenExpInHours,
		RefreshTokenMaxAge: refreshTokenExpInDays,
	}

	return config, nil
}
