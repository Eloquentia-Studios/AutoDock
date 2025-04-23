package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type config struct {
	Interval time.Duration `env:"INTERVAL,required" envDefault:"30m"`
}

var (
	Interval time.Duration
)

func Load() {
	_ = godotenv.Load()

	var cfg config
	err := env.Parse(&cfg)

	if err != nil {
		panic(err) // Handle error
	}

	Interval = cfg.Interval
}
