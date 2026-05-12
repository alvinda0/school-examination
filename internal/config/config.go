package config

import "os"

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	ServerPort  string
	JWTSecret   string
	CORSOrigins string
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5433"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "123456"),
		DBName:      getEnv("DB_NAME", "postgres"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
