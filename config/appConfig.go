package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort string
	Dsn        string
	AppSecret  string

	BrevoAPIKey string
	EmailFrom   string

	StripeSecret string
	SuccessUrl   string
	CancelUrl    string
}

func SetupEnvironment() (AppConfig, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env not loaded")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, errors.New("HTTP_PORT not set")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return AppConfig{}, errors.New("DSN not set")
	}

	appSecret := os.Getenv("APP_SECRET")
	if appSecret == "" {
		return AppConfig{}, errors.New("APP_SECRET not set")
	}

	brevoAPIKey := os.Getenv("BREVO_API_KEY")
	if brevoAPIKey == "" {
		return AppConfig{}, errors.New("BREVO_API_KEY not set")
	}

	emailFrom := os.Getenv("EMAIL_FROM")
	if emailFrom == "" {
		return AppConfig{}, errors.New("EMAIL_FROM not set")
	}

	return AppConfig{
		ServerPort:   httpPort,
		Dsn:          dsn,
		AppSecret:    appSecret,
		BrevoAPIKey:  brevoAPIKey,
		EmailFrom:    emailFrom,
		StripeSecret: os.Getenv("STRIPE_SECRET"),
		SuccessUrl:   os.Getenv("SUCCESS_URL"),
		CancelUrl:    os.Getenv("CANCEL_URL"),
	}, nil
}
