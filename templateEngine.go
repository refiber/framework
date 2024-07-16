package refiber

import (
	"html/template"

	"github.com/gofiber/template/html/v2"

	"github.com/refiber/framework/vite"
)

func newTemplateEngine() *html.Engine {
	engine := html.New("./resources/views", ".tpl")

	engine.AddFunc(
		"raw", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	engine.AddFunc("vite", func(s ...string) template.HTML {
		return template.HTML(vite.GetScripts(s...))
	})

	engine.AddFunc("inertia", func(s string) template.HTML {
		return template.HTML(s)
	})

	return engine
}
