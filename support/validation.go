package support

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/constant"
	"github.com/refiber/framework/util"
)

// TODO: support multi lang

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

	var nestedErrors []fiber.Map
	errors := make(fiber.Map, len(err.(validator.ValidationErrors)))
	for _, err := range err.(validator.ValidationErrors) {
		if parts := strings.Split(err.Namespace(), "."); len(parts) > 2 {
			nestedError := v.getNestedErrorByNamespaceParts(parts, err)
			nestedErrors = append(nestedErrors, *nestedError)
			continue
		}

		errors[err.Field()] = err.Translate(*v.support.translator)
	}

	if len(nestedErrors) > 0 {
		mergedNestedErrors := *util.DeepMergeFiberMaps(nestedErrors...)
		errors = *util.MergeFiberMaps(errors, mergedNestedErrors)
	}

	if err := v.sharedData.saveTempData(constant.SessionKeyError, &errors); err != nil {
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

func (v *validation) removeBrackets(input string) string {
	index := strings.Index(input, "[")
	if index != -1 {
		return input[:index]
	}
	return input
}

func (v *validation) getBracketValue(input string) *string {
	start := strings.Index(input, "[")
	end := strings.Index(input, "]")

	var value *string

	if start != -1 && end != -1 && start < end {
		v := input[start+1 : end]
		value = &v
	}

	return value
}

func (v *validation) getNestedErrorByNamespaceParts(parts []string, err validator.FieldError) *fiber.Map {
	var (
		nedtedError, prev *fiber.Map
		prevKey           *string
	)

	for i, key := range parts {
		if i == 0 {
			continue
		}

		var indexKey *string
		if strings.Contains(parts[i-1], "[") {
			indexKey = v.getBracketValue(parts[i-1])
		}

		if strings.Contains(key, "[") {
			key = v.removeBrackets(key)
		}

		if prevKey == nil {
			prevKey = &key
		}

		if prev == nil {
			data := fiber.Map{*prevKey: fiber.Map{}}
			prev = &data
			nedtedError = &data
		} else {
			data := (*prev)[*prevKey].(fiber.Map)

			if err.Field() == key {
				v := err.Translate(*v.support.translator)
				if indexKey != nil {
					data[*indexKey] = fiber.Map{key: v}
				} else {
					data[key] = v
				}
			} else {
				newData := data[*prevKey].(fiber.Map)

				if indexKey != nil {
					newData[*indexKey] = fiber.Map{key: fiber.Map{}}
				} else {
					newData[key] = fiber.Map{}
				}
				prev = &newData
			}
		}

		prevKey = &key
	}

	return nedtedError
}
