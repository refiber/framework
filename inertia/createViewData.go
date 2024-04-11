package inertia

import (
	"encoding/json"
	"fmt"
	"html"
	"os"

	"github.com/gofiber/fiber/v2"
)

func createViewData(bufJsonProps *[]byte, viewData *fiber.Map) *fiber.Map {
	data := fiber.Map{}
	if viewData != nil {
		data = *viewData
	}

	data["refiber"] = ""

	if buf, err := os.ReadFile("./public/hot"); err == nil {
		viteURL := string(buf)
		scripts := fmt.Sprintf(`
			<script type="module">
				import RefreshRuntime from "%s/@react-refresh"
				RefreshRuntime.injectIntoGlobalHook(window)
				window.$RefreshReg$ = () => {}
				window.$RefreshSig$ = () => (type) => type
				window.__vite_plugin_react_preamble_installed__ = true
			</script>
			<script type="module" src="%s/@vite/client"></script>
			<script type="module" src="%s/resources/js/app.tsx"></script>
		`, viteURL, viteURL, viteURL)

		data["refiber"] = scripts
	} else if buf, err := os.ReadFile("./public/build/manifest.json"); err == nil {
		manifest := make(map[string]interface{})
		if err := json.Unmarshal(buf, &manifest); err == nil {
			main := manifest["resources/js/app.tsx"].(map[string]interface{})
			js := main["file"].(string)
			css := main["css"].([]interface{})

			cssScripts := ""
			for _, c := range css {
				cssScripts += fmt.Sprintf(`
					<link rel="stylesheet" href="/build/%s" />
				`, c)
			}

			scripts := fmt.Sprintf(`
				<script type="module" src="/build/%s"></script>	
				%s
			`, js, cssScripts)

			data["refiber"] = scripts
		}
	}

	data["refiber"] = fmt.Sprintf(`
		%s
		<div id="app" data-page=%#v></div> 
	`, data["refiber"], html.EscapeString(string(*bufJsonProps)))

	return &data
}
