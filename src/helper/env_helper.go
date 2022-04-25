package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	env map[string]string
)

func SetEnv() {
	env_err := godotenv.Load(".env")
	if env_err != nil {
		log.Fatalf("failed to load env file: %s", env_err)
	}
	env = make(map[string]string)
	env["PORT"] = os.Getenv("PORT")
	env["PG_HOST"] = os.Getenv("PG_HOST")
	env["PG_PORT"] = os.Getenv("PG_PORT")
	env["PG_USER"] = os.Getenv("PG_USER")
	env["PG_PASS"] = os.Getenv("PG_PASS")
	env["PG_DB"] = os.Getenv("PG_DB")
}

func GetEnv(key string) string {
	return env[key]
}
