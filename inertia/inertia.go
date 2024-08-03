package inertia

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/support"
	"github.com/refiber/framework/utils"
	"github.com/refiber/framework/vite"
)

type InertiaInterface interface {
	SetViewTemplate(view string)
	Render() *render
}

type PreRenderHanlder = func(string) *string

type Config struct {
	App                      support.Refiber
	PreRenderHanlder         PreRenderHanlder
	EnablePreRenderByDefault bool
	ViewTemplate             string
}

func New(c Config) *inertia {
	i := inertia{
		s:                        c.App,
		viewTemplate:             "app",
		PreRenderHanlder:         c.PreRenderHanlder,
		EnablePreRenderByDefault: c.EnablePreRenderByDefault,
	}

	if c.ViewTemplate != "" {
		i.viewTemplate = c.ViewTemplate
	}

	return &i
}

type inertia struct {
	s                        support.Refiber
	viewTemplate             string
	PreRenderHanlder         PreRenderHanlder
	EnablePreRenderByDefault bool
}

func (i *inertia) Render() *render {
	r := &render{}
	r.inertia = i
	r.viewTemplate = i.viewTemplate
	r.preRender = i.EnablePreRenderByDefault

	return r
}

func (i *inertia) SetViewTemplate(view string) {
	i.viewTemplate = view
}

type render struct {
	*inertia
	viewData  *fiber.Map
	preRender bool
}

func (r *render) SetViewData(data *fiber.Map) *render {
	r.viewData = data
	return r
}

func (r *render) DisablePreRender() *render {
	r.preRender = false
	return r
}

func (r *render) Page(page string, props *fiber.Map) error {
	sharedProps := fiber.Map{}

	if session, err := r.s.GetSession().Get(r.s.GetCtx()); err == nil {
		sharedProps = *support.GetTempData(session)
	}

	data := fiber.Map{}
	data["url"] = r.s.GetCtx().OriginalURL()
	v := utils.GetMD5Hash("./public/build/manifest.json")
	data["version"] = v
	data["component"] = page
	data["props"] = utils.MergeFiberMaps(&sharedProps, props)

	headers := r.s.GetCtx().GetReqHeaders()

	headerXInertia, exist := headers["X-Inertia"]
	headerXInertiaVersion, exist2 := headers["X-Inertia-Version"]
	renderViewTemplate := !exist || !exist2 || len(headerXInertia) > 0 && headerXInertia[0] != "true" || len(headerXInertiaVersion) > 0 && headerXInertiaVersion[0] != v

	if renderViewTemplate {
		jsonProps, _ := json.Marshal(data)
		viewData := createViewData(&jsonProps, r.viewData)

		err := r.s.GetCtx().Render(r.viewTemplate, viewData)

		// pre-render only available on production
		viteDevURL := vite.GetDevelopmentURL()
		if viteDevURL == nil && r.PreRenderHanlder != nil && r.preRender {
			if manifest := vite.GetManifest(); manifest != nil {
				// file = "assets/app*.js"
				if file := manifest.GetFileByResource(rootAppFile); file != nil {
					if scriptBuf, err := os.ReadFile(fmt.Sprintf(`./public/build/%s`, *file)); err == nil {
						html := string(r.s.GetCtx().Response().Body())

						oldScript := vite.CreateScriptTag(fmt.Sprintf(`/build/%s`, *file))
						newScript := fmt.Sprintf(`<script type="module">%s</script>`, string(scriptBuf))

						html = strings.Replace(html, oldScript, newScript, 1)
						preRenderedHtml := r.PreRenderHanlder(html)

						if preRenderedHtml != nil {
							html = strings.Replace(*preRenderedHtml, newScript, oldScript, 1)
							r.s.GetCtx().Response().SetBody([]byte(html))
						}
					}
				}
			}
		}

		return err
	}

	r.s.GetCtx().Response().Header.Set("X-Inertia", "true")

	return r.s.GetCtx().Status(fiber.StatusOK).JSON(data)
}
