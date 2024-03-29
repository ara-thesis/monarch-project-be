package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	env map[string]string = make(map[string]string)
)

type EnvHelper struct {
	EnvInterface
}

type EnvInterface interface {
	SetEnv()
	GetEnv()
}

func SetEnv() {
	env_err := godotenv.Load(".env")
	if env_err != nil {
		log.Printf("failed to load env file: %s", env_err)
		return
	}
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
