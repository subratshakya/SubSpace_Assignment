package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	LinkedInEmail    string
	LinkedInPassword string
	Headless         bool
	MaxActions       int
	ConnectionNote   string
	SearchURL        string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	
	config := &Config{
		LinkedInEmail:    getEnv("LINKEDIN_EMAIL", ""),
		LinkedInPassword: getEnv("LINKEDIN_PASSWORD", ""),
		Headless:         getEnvBool("HEADLESS", false),
		MaxActions:       getEnvInt("MAX_ACTIONS", 5),
		ConnectionNote:   getEnv("CONNECTION_NOTE", "Hi! I'd like to connect with you."),
		SearchURL:        getEnv("SEARCH_URL", ""),
	}
	
	// Validate required fields
	if config.LinkedInEmail == "" {
		return nil, fmt.Errorf("LINKEDIN_EMAIL is required")
	}
	if config.LinkedInPassword == "" {
		return nil, fmt.Errorf("LINKEDIN_PASSWORD is required")
	}
	
	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool gets a boolean environment variable
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return boolValue
}

// getEnvInt gets an integer environment variable
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}

