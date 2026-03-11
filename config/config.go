package config

import (
	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	FolderId   string `env:"FOLDER_ID,required"`
	ApiKey     string `env:"API_KEY,required"`
	NatsUrl    string `env:"NATS_URL,required"`
	KafkaAddr  string `env:"KAFKA_ADDR,required"`
	KafkaTopic string `env:"KAFKA_TOPIC,required"`
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("env variables NOT loaded from .env file")
		return nil, err
	}
	cfg := new(Config)
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
