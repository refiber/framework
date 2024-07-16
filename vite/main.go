package vite

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// TODO: find a way to call this only once every request
func GetDevelopmentURL() *string {
	// TODO: skip if env == prod
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
      `, getScirptTag(fmt.Sprintf("%s%s", *developmentURL, r)))
		}

		scripts := fmt.Sprintf(`
			<script type="module" src="%s/@vite/client"></script>
      %s
		`, *developmentURL, data)

		data = scripts
	} else if buf, err := os.ReadFile("./public/build/manifest.json"); err == nil {
		data = getCompailedScripts(resources, buf)
	}

	return data
}

func getCompailedScripts(resources []string, buf []byte) string {
	var data string

	manifest := make(map[string]interface{})
	err := json.Unmarshal(buf, &manifest)
	if err != nil {
		return data
	}

	for _, _r := range resources {
		r := _r
		// remove / on first char
		if string(r[0]) == "/" {
			r = string(r[1:])
		}

		main := manifest[r].(map[string]interface{})
		file, ok := main["file"].(string)
		if !ok {
			continue
		}

		scripts := getScirptTag(fmt.Sprintf("/build/%s", file))

		if css, ok := main["css"].([]interface{}); ok {
			cssScripts := ""

			for _, c := range css {
				cssScripts += fmt.Sprintf(`
        %s
      `, getScirptTag(fmt.Sprintf("/build/%s", c)))
			}

			scripts += fmt.Sprintf(`
        %s
      `, cssScripts)
		}

		data += scripts
	}
	return data
}

func getScirptTag(source string) string {
	if strings.Contains(source, ".css") {
		return fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, source)
	}
	return fmt.Sprintf(`<script type="module" src="%s"></script>`, source)
}
