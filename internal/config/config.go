package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	BackEndPort int    `yaml:"port" env:"BACKEND_PORT" env-default:"8080"`
}

func InitConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
		log.Fatalf("can`t read config: %s", err)
	}
	return &cfg
}
