package refiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/badger/v2"
)

type Config struct {
	AppName        string
	SessionStorage fiber.Storage
}

func configDefault(config ...Config) Config {
	cfg := Config{}

	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.AppName == "" {
		cfg.AppName = "Refiber"
	}

	if cfg.SessionStorage == nil {
		storage := badger.New(badger.Config{
			Database:   "./storage/framework/session.badger",
			Reset:      false,
			GCInterval: 10 * time.Second,
		})
		cfg.SessionStorage = storage
	}

	return cfg
}
