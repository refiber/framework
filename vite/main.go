package vite

import (
	"fmt"
	"os"
	"strings"
)

// TODO: find a way to call this only once every request
func GetDevelopmentURL() *string {
	env := os.Getenv("APP_ENV")
	if env != "local" && env != "dev" {
		return nil
	}

	if buf, err := os.ReadFile("./public/hot"); err == nil {
		url := string(buf)
		return &url
	}

	return nil
}

func GetScripts(resources ...string) string {
	var data string

	developmentURL := GetDevelopmentURL()

	if developmentURL != nil {
		for _, _r := range resources {
			r := _r
			// add / on first char
			if string(_r[0]) != "/" {
				r = "/" + _r
			}

			data += fmt.Sprintf(`
        %s
      `, CreateScriptTag(fmt.Sprintf("%s%s", *developmentURL, r)))
		}

		scripts := fmt.Sprintf(`
			<script type="module" src="%s/@vite/client"></script>
      %s
		`, *developmentURL, data)

		data = scripts
	} else if m := GetManifest(); m.data != nil {
		data = getCompailedScripts(resources, *m)
	}

	return data
}

func getCompailedScripts(resources []string, m manifest) string {
	var data string

	for _, r := range resources {
		file := m.GetFileByResource(r)
		if file == nil {
			continue
		}

		scripts := CreateScriptTag(fmt.Sprintf("/build/%s", *file))

		if css := m.GetCSSbyResource(r); css != nil {
			cssScripts := ""

			for _, c := range *css {
				cssScripts += fmt.Sprintf(`
        %s
      `, CreateScriptTag(fmt.Sprintf("/build/%s", c)))
			}

			scripts += fmt.Sprintf(`
        %s
      `, cssScripts)
		}

		data += scripts
	}
	return data
}

func CreateScriptTag(source string) string {
	if strings.Contains(source, ".css") {
		return fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, source)
	}
	return fmt.Sprintf(`<script type="module" src="%s"></script>`, source)
}
