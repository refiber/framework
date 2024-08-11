package support

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Refiber interface {
	GetSessionStore() *session.Store
	GetValidator() *validator.Validate
	GetTranslator() *ut.Translator

	Render(*fiber.Ctx) *render
	Redirect(*fiber.Ctx) *redirect
	Auth(*fiber.Ctx) *auth
	SharedData(*fiber.Ctx) *sharedData
	Validation(*fiber.Ctx) *validation
}

func NewSupport(session *session.Store, validate *validator.Validate, translator *ut.Translator) *support {
	return &support{session, validate, translator}
}

type support struct {
	sessionStore *session.Store
	validator    *validator.Validate
	translator   *ut.Translator
}

func (s *support) GetSessionStore() *session.Store {
	return s.sessionStore
}

func (s *support) GetValidator() *validator.Validate {
	return s.validator
}

func (s *support) GetTranslator() *ut.Translator {
	return s.translator
}
