package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/psyb0t/gonfiguration"
)

type AppConfig struct {
	// Basic types
	ListenAddress string `env:"LISTEN_ADDRESS"`
	Debug         bool   `env:"DEBUG"`
	Port          int    `env:"PORT"`
	
	// Advanced types
	Timeout       time.Duration `env:"TIMEOUT"`
	AllowedHosts  []string      `env:"ALLOWED_HOSTS"`
	
	// Database shit
	DBDSN    string `env:"DB_DSN"`
	DBName   string `env:"DB_NAME"`
	DBUser   string `env:"DB_USER"`
	DBPass   string `env:"DB_PASS"`
}

func main() {
	cfg := AppConfig{}

	// Set some defaults because you're not a savage
	gonfiguration.SetDefaults(map[string]interface{}{
		"LISTEN_ADDRESS": "127.0.0.1:8080",
		"TIMEOUT":        30 * time.Second,
		"ALLOWED_HOSTS":  []string{"localhost", "127.0.0.1"},
		"DB_DSN":         "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	})

	// Set some env vars (in real life these come from your environment)
	os.Setenv("PORT", "8080")
	os.Setenv("DEBUG", "true")
	os.Setenv("DB_NAME", "myapp")
	os.Setenv("DB_USER", "postgres-user")
	os.Setenv("DB_PASS", "super-secret-password")
	os.Setenv("ALLOWED_HOSTS", "api.example.com, cdn.example.com, *.example.com")

	// Parse that shit
	if err := gonfiguration.Parse(&cfg); err != nil {
		log.Fatalf("holy fuque! config parsing failed: %v", err)
	}

	fmt.Printf("Config loaded: %+v\n", cfg)
	fmt.Printf("Allowed hosts: %v\n", cfg.AllowedHosts) // ["api.example.com", "cdn.example.com", "*.example.com"]
	fmt.Printf("Timeout: %v\n", cfg.Timeout)           // 30s
}
