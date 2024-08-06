package vite

import (
	"encoding/json"
	"os"
)

type ManifestInterface interface {
	GetDataByResource(resource string) map[string]interface{}
	GetCompailedFileNameByResource(resource string) *string
	GetCompailedCSSFileNamesByResource(resource string) *[]*string
	GetCompailedImportFileNamesByResource(resource string) *[]*string
}

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

// will return "assets/app*.js"
func (m *manifest) GetCompailedFileNameByResource(resource string) *string {
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

func (m *manifest) GetCompailedCSSFileNamesByResource(resource string) *[]*string {
	return m.getManyFileNameByResourceAndKey(resource, "css")
}

func (m *manifest) GetCompailedImportFileNamesByResource(resource string) *[]*string {
	return m.getManyFileNameByResourceAndKey(resource, "imports")
}

func (m *manifest) getManyFileNameByResourceAndKey(resource, key string) *[]*string {
	main := m.GetDataByResource(resource)
	if main == nil {
		return nil
	}

	data, ok := main[key].([]interface{})
	if !ok {
		return nil
	}

	var results []*string

	for _, v := range data {
		value, ok := v.(string)
		if !ok {
			continue
		}

		results = append(results, &value)
	}

	return &results
}
