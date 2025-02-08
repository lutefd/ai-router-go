package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         int
	OPENAI_SK          string
	DEEPSEEK_SK        string
	GEMINI_SK          string
	MongoDBURI         string
	MongoDBDatabase    string
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
	ClientURL          string
	AuthRedirectURL    string
	AndroidClientID    string
	WorkerID           int64
}

func LoadConfig(skipEnvFile ...bool) (*Config, error) {
	if len(skipEnvFile) == 0 || !skipEnvFile[0] {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file, falling back to system environment variables")
		}
	}

	config := &Config{}

	srvPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("SERVER_PORT environment variable is not a valid integer")
	}
	config.ServerPort = srvPort

	config.OPENAI_SK = os.Getenv("OPENAI_SK")
	if config.OPENAI_SK == "" {
		return nil, fmt.Errorf("OPENAI_SK environment variable is not set")
	}

	config.DEEPSEEK_SK = os.Getenv("DEEPSEEK_SK")
	if config.DEEPSEEK_SK == "" {
		return nil, fmt.Errorf("DEEPSEEK_SK environment variable is not set")
	}

	config.GEMINI_SK = os.Getenv("GEMINI_SK")
	if config.GEMINI_SK == "" {
		return nil, fmt.Errorf("GEMINI_SK environment variable is not set")
	}

	config.MongoDBURI = os.Getenv("MONGODB_URI")
	if config.MongoDBURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	config.MongoDBDatabase = os.Getenv("MONGODB_DATABASE")
	if config.MongoDBDatabase == "" {
		return nil, fmt.Errorf("MONGODB_DATABASE environment variable is not set")
	}

	config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if config.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable is not set")
	}

	config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	if config.GoogleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET environment variable is not set")
	}

	config.JWTSecret = os.Getenv("JWT_SECRET")
	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	config.ClientURL = os.Getenv("CLIENT_URL")
	if config.ClientURL == "" {
		return nil, fmt.Errorf("CLIENT_URL environment variable is not set")
	}

	config.AuthRedirectURL = os.Getenv("AUTH_REDIRECT_URL")
	if config.AuthRedirectURL == "" {
		return nil, fmt.Errorf("AUTH_REDIRECT_URL environment variable is not set")
	}

	config.WorkerID, err = strconv.ParseInt(os.Getenv("WORKER_ID"), 10, 64)
	if err != nil || config.WorkerID < 1 {
		config.WorkerID = int64(1)
		log.Println("WORKER_ID environment variable is not set or invalid, defaulting to 1")
	}
	config.AndroidClientID = os.Getenv("ANDROID_CLIENT_ID")
	if config.AndroidClientID == "" {
		return nil, fmt.Errorf("ANDROID_CLIENT_ID environment variable is not set")
	}
	return config, nil
}
