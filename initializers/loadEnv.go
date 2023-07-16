package initializers

import (
	"log"
	"os"
	"strconv"
	"time"

	// go dotenv
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	error := godotenv.Load()
	if error != nil {
		log.Fatal("Error loading .env file")
	}

	accessTokenExpInHours, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRED_IN_HOURS"))
	refreshTokenExpInDays, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRED_IN_DAYS"))

	config = Config{
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
