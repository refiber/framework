package support

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/constant"
)

func (s *support) Validate(sct interface{}) error {
	err := s.validate.Struct(sct)
	if err == nil {
		return nil
	}

	m := make(fiber.Map, len(err.(validator.ValidationErrors)))
	for _, err := range err.(validator.ValidationErrors) {
		// use lower case?
		m[strings.ToLower(err.Field())] = err.Translate(*s.translator)
	}

	if err := saveTempData(s, constant.SessionKeyError, &m); err != nil {
		return err
	}

	return err
}

type ValidationErrorField struct {
	Name    string
	Message string
}

func (s *support) CreateValidationErrors(fields []*ValidationErrorField) error {
	session, _ := s.session.Get(s.GetCtx())
	sessionKey := string(constant.SessionKeyError) + session.ID()

	m := fiber.Map{}
	for _, f := range fields {
		m[f.Name] = f.Message
	}

	buf, _ := json.Marshal(m)
	session.Set(sessionKey, buf)
	session.SetExpiry(time.Minute * 1)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}
