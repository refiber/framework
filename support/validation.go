package support

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/constant"
)

func (s *support) Validation(ctx *fiber.Ctx) *validation {
	return &validation{s, ctx, s.SharedData(ctx)}
}

type validation struct {
	support    *support
	ctx        *fiber.Ctx
	sharedData *sharedData
}

func (v *validation) Validate(sct interface{}) error {
	err := v.support.validator.Struct(sct)
	if err == nil {
		return nil
	}

	m := make(fiber.Map, len(err.(validator.ValidationErrors)))
	for _, err := range err.(validator.ValidationErrors) {
		m[err.Field()] = err.Translate(*v.support.translator)
	}

	if err := v.sharedData.saveTempData(constant.SessionKeyError, &m); err != nil {
		return err
	}

	return err
}

type ValidationErrorField struct {
	Name    string
	Message string
}

func (v *validation) SetErrors(fields []*ValidationErrorField) error {
	m := fiber.Map{}
	for _, f := range fields {
		m[f.Name] = f.Message
	}

	if err := v.sharedData.saveTempData(constant.SessionKeyError, &m); err != nil {
		return err
	}

	return nil
}
