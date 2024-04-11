package support

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

func saveTempSession(s *support, key string, data *fiber.Map) error {
	session, _ := s.session.Get(s.GetCtx())

	buf, _ := json.Marshal(data)
	session.Set(key+session.ID(), buf)
	session.SetExpiry(time.Minute * 1)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}
