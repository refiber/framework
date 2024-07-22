package support

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Refiber interface {
	Render() *render
	Validate(s interface{}) error
	CreateValidationErrors(fields []*ValidationErrorField) error
	NewAuthenticatedUserSession(user interface{}) error
	UpdateAuthenticatedUserSession(user interface{}) error
	GetAuthenticatedUserSession(user interface{}) error
	DestroyAuthenticatedUserSession() error
	GetSession() *session.Store
	GetValidate() *validator.Validate
	GetTranslator() *ut.Translator
	Redirect() *redirect
	GetCtx() *fiber.Ctx
	SetSharedData(data fiber.Map) error
}

func NewSupport(session *session.Store, validate *validator.Validate, translator *ut.Translator) *support {
	return &support{session, validate, translator, nil}
}

type support struct {
	session    *session.Store
	validate   *validator.Validate
	translator *ut.Translator
	Ctx        *fiber.Ctx
}

func (s *support) GetSession() *session.Store {
	return s.session
}

func (s *support) GetValidate() *validator.Validate {
	return s.validate
}

func (s *support) GetTranslator() *ut.Translator {
	return s.translator
}

func (s *support) GetCtx() *fiber.Ctx {
	return s.Ctx
}

func (s *support) SetCtx(ctx *fiber.Ctx) {
	s.Ctx = ctx
}
