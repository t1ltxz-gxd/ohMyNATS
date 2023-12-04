package misc

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadDotEnv(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
