package config

import "os"

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      port,
	}
}