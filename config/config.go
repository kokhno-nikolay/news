package config

import (
	"encoding/json"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var (
	once sync.Once
	cfg  *Config
)

type Config struct {
	GrpcAddress string `env:"SERVER_PORT" envDefault:"localhost:50051"`
	HttpAddress string `env:"HTTP_ADDRESS" envDefault:"localhost:8000"`
	PostresDNS  string `env:"POSTGRES_DNS"`
}

func (c *Config) String() string {
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}

func GetConfig() *Config {
	_ = godotenv.Load()

	once.Do(func() {
		config := Config{}
		if err := env.Parse(&config); err != nil {
			panic(err)
		}

		cfg = &config
	})

	return cfg
}
