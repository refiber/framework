package inertia

import (
	"fmt"
	"html"

	"github.com/gofiber/fiber/v2"
)

var rootAppFile = "/resources/js/app.tsx"

func createViewData(bufJsonProps *[]byte, viewData *fiber.Map) *fiber.Map {
	data := fiber.Map{}
	if viewData != nil {
		data = *viewData
	}

	data["props"] = fmt.Sprintf(`<div id="app" data-page=%#v></div>`, html.EscapeString(string(*bufJsonProps)))

	return &data
}
