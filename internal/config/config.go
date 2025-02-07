package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  int
	OPENAI_SK   string
	DEEPSEEK_SK string
	GEMINI_SK   string
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

	return config, nil
}
