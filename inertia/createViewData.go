package inertia

import (
	"fmt"
	"html"

	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/vite"
)

// TODO: make this customizable?
var rootAppFile = "/resources/js/app.tsx"

type viewDataStruct struct {
	ScriptTags  string
	DataPageDiv string
	AllTags     string
}

func createViewData(bufJsonProps *[]byte, viewData *fiber.Map) (*fiber.Map, *viewDataStruct) {
	data := fiber.Map{}
	if viewData != nil {
		data = *viewData
	}

	data["props"] = vite.GetScripts(rootAppFile)

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

	vds := viewDataStruct{
		ScriptTags:  data["props"].(string),
		DataPageDiv: fmt.Sprintf(`<div id="app" data-page=%#v></div>`, html.EscapeString(string(*bufJsonProps))),
	}

	vds.AllTags = fmt.Sprintf(`
		%s
		%s
	`, vds.ScriptTags, vds.DataPageDiv)

	data["props"] = vds.AllTags

	return &data, &vds
}
