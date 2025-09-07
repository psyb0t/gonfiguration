package main

import (
	"fmt"
	"log"
	"os"

	"github.com/psyb0t/gonfiguration"
)

type config struct {
	ListenAddress string `env:"LISTEN_ADDRESS"`
	DBDSN         string `env:"DB_DSN"`
	DBName        string `env:"DB_NAME"`
	DBUser        string `env:"DB_USER"`
	DBPass        string `env:"DB_PASS"`
}

func main() {
	cfg := config{}

	gonfiguration.SetDefaults(map[string]interface{}{
		"LISTEN_ADDRESS": "127.0.0.1:8080",
		"DB_DSN":         "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	})

	if err := os.Setenv("DB_NAME", "postgres"); err != nil {
		log.Fatalf("holy fuque! can't set env: %v", err)
	}

	if err := os.Setenv("DB_USER", "postgres-user"); err != nil {
		log.Fatalf("holy fuque! can't set env: %v", err)
	}

	if err := os.Setenv("DB_PASS", "postgres-pass"); err != nil {
		log.Fatalf("holy fuque! can't set env: %v", err)
	}

	if err := gonfiguration.Parse(&cfg); err != nil {
		log.Fatalf("holy fuque! can't parse config: %v", err)
	}

	fmt.Printf("%+v\n", cfg) //nolint:forbidigo
}
