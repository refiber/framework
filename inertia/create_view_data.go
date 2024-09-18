package inertia

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var rootAppFile = "/resources/js/app.tsx"

func createViewData(bufJsonProps *[]byte, viewData *fiber.Map) *fiber.Map {
	data := fiber.Map{}
	if viewData != nil {
		data = *viewData
	}

	data["props"] = fmt.Sprintf(`<div id="app" data-page=%s></div>`, string(*bufJsonProps))

	return &data
}
