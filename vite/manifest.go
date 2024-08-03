package vite

import (
	"encoding/json"
	"os"
)

type manifest struct {
	data map[string]interface{}
}

func GetManifest() *manifest {
	data := make(map[string]interface{})

	buf, err := os.ReadFile("./public/build/manifest.json")
	if err != nil {
		return nil
	}

	if err = json.Unmarshal(buf, &data); err != nil {
		return nil
	}

	return &manifest{data: data}
}

// will return value of resource key, example: "resources/js/app.tsx": {...}
func (m *manifest) GetDataByResource(resource string) map[string]interface{} {
	r := resource
	// remove / on first char
	if string(r[0]) == "/" {
		r = string(r[1:])
	}

	if m.data == nil {
		return nil
	}

	main, ok := m.data[r].(map[string]interface{})
	if !ok {
		return nil
	}

	return main
}

// will return value of file key, example: "file": "assets/app.js"
func (m *manifest) GetFileByResource(resource string) *string {
	main := m.GetDataByResource(resource)
	if main == nil {
		return nil
	}

	file, ok := main["file"].(string)
	if !ok {
		return nil
	}

	return &file
}

func (m *manifest) GetCSSbyResource(resource string) *[]interface{} {
	main := m.GetDataByResource(resource)
	if main == nil {
		return nil
	}

	css, ok := main["css"].([]interface{})
	if !ok {
		return nil
	}

	return &css
}
