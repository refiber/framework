package refiber

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	AppName        string
	SessionStorage fiber.Storage
}

func configDefault(cfg Config) Config {
	if cfg.SessionStorage == nil {
		panic("[config]: SessionStorage can't be null")
	}

	if cfg.AppName == "" {
		cfg.AppName = "Refiber"
	}

	return cfg
}
