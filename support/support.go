package support

import (
	"encoding/json"
	"refiber/constant"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Refiber interface {
	Render() *render
	Validate(s interface{}) error
	CreateValidationErrors(fields []*ValidationErrorField) error
	NewAuthenticatedUserSession(user interface{}) error
	GetAuthenticatedUserSession(user interface{}) error
	DestroyAuthenticatedUserSession() error
	GetSession() *session.Store
	GetValidate() *validator.Validate
	GetTranslator() *ut.Translator
  Redirect() *redirect
	GetCtx() *fiber.Ctx
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

// TODO: make this customizeable
func GetSharedMap(session *session.Session) *fiber.Map {
	m := make(fiber.Map)
	m["errors"] = fiber.Map{}
	m["auth"] = new(fiber.Map)
	m["flash"] = new(fiber.Map)

	/**
	 * Form Errors
	 */
	keyErrors := constant.KeyErrors + session.ID()
	raw := session.Get(keyErrors)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetSharedMap: failed to get errors")
		} else {
			m["errors"] = d
			session.Delete(keyErrors)
		}
	}

	/**
	 * Flash Message
	 */
	keyFlashMessage := constant.KeyFlashMessage + session.ID()
	raw = session.Get(keyFlashMessage)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetSharedMap: failed to get keyFlashMessage")
		} else {
			m["flash"] = d
			session.Delete(keyFlashMessage)
		}
	}

	/**
	 * Auth
	 */
	keyAuth := constant.KeyAuth + session.ID()
	raw = session.Get(keyAuth)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetSharedMap: failed to get auth")
		} else {
			m["auth"] = d
		}
	}

	if err := session.Save(); err != nil {
		log.Errorw("refiber.support.GetSharedMap: failed to save session")
	}

	return &m
}
