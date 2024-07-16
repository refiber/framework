package inertia

import (
	"fmt"
	"html"

	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/vite"
)

func createViewData(bufJsonProps *[]byte, viewData *fiber.Map) *fiber.Map {
	data := fiber.Map{}
	if viewData != nil {
		data = *viewData
	}

	data["props"] = vite.GetScripts("/resources/js/app.tsx")

	viteDevelopmentURL := vite.GetDevelopmentURL()

	// TODO: also check if app using react
	if viteDevelopmentURL != nil {
		data["props"] = fmt.Sprintf(`
    <script type="module">
      import RefreshRuntime from "%s/@react-refresh"
      RefreshRuntime.injectIntoGlobalHook(window)
      window.$RefreshReg$ = () => {}
      window.$RefreshSig$ = () => (type) => type
      window.__vite_plugin_react_preamble_installed__ = true
    </script>
    %s
  `, *viteDevelopmentURL, data["props"])
	}

	data["props"] = fmt.Sprintf(`
		%s
		<div id="app" data-page=%#v></div> 
	`, data["props"], html.EscapeString(string(*bufJsonProps)))

	return &data
}
