package main

import (
	"github.com/joho/godotenv"
	"os"
)

func loadEnv() {
	env := os.Getenv("MADDERMIND_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}

	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}
