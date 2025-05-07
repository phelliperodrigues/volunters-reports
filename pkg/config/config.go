package config

import (
	"os"
	"path/filepath"
)

// Config holds all configuration values
type Config struct {
	// File paths
	BooksFilePath    string
	InputFilePath    string
	OutputDirPath    string
	ReportsDirPath   string
	PDFOutputDirPath string

	// CSV configuration
	CSVHeaderRows int

	// Server configuration
	ServerPort string
}

// New creates a new configuration with default values
func New() *Config {
	basePath := getBasePath()

	return &Config{
		BooksFilePath:    filepath.Join(basePath, "files", "books.csv"),
		InputFilePath:    filepath.Join(basePath, "files", "input.csv"),
		OutputDirPath:    filepath.Join(basePath, "files", "output"),
		ReportsDirPath:   filepath.Join(basePath, "data", "reports"),
		PDFOutputDirPath: filepath.Join(basePath, "files", "output"),
		CSVHeaderRows:    12,
		ServerPort:       getEnvOrDefault("SERVER_PORT", "8080"),
	}
}

// getBasePath returns the base path of the project
func getBasePath() string {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "."
	}
	return basePath
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
