package inertia

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/support"
	"github.com/refiber/framework/utils"
	"github.com/refiber/framework/vite"
)

type InertiaInterface interface {
	SetViewTemplate(view string)
	Render() *render
}

type (
	PreRenderHanlder = func(PreRenderInterface) *string
	SSRHanlder       = func(SSRInterface) *string
)

type Config struct {
	App                      support.Refiber
	PreRenderHanlder         PreRenderHanlder
	SSRHanlder               SSRHanlder
	EnablePreRenderByDefault bool
	EnableSSRByDefault       bool
	ViewTemplate             string
}

func New(c Config) *inertia {
	i := inertia{
		s:                        c.App,
		viewTemplate:             "app",
		PreRenderHanlder:         c.PreRenderHanlder,
		SSRHanlder:               c.SSRHanlder,
		EnablePreRenderByDefault: c.EnablePreRenderByDefault,
		EnableSSRByDefault:       c.EnableSSRByDefault,
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
	SSRHanlder               SSRHanlder
	EnablePreRenderByDefault bool
	EnableSSRByDefault       bool
}

func (i *inertia) Render() *render {
	r := &render{}
	r.inertia = i
	r.viewTemplate = i.viewTemplate
	r.preRender = i.EnablePreRenderByDefault
	r.ssr = i.EnableSSRByDefault

	return r
}

func (i *inertia) SetViewTemplate(view string) {
	i.viewTemplate = view
}

type render struct {
	*inertia
	viewData  *fiber.Map
	preRender bool
	ssr       bool
}

func (r *render) SetViewData(data *fiber.Map) *render {
	r.viewData = data
	return r
}

func (r *render) DisablePreRender() *render {
	r.preRender = false
	return r
}

func (r *render) EnablePreRender() *render {
	r.preRender = true
	return r
}

func (r *render) DisableSSR() *render {
	r.ssr = false
	return r
}

func (r *render) EnableSSR() *render {
	r.ssr = true
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
		viewData, viewDataStruct := createViewData(&jsonProps, r.viewData)

		err := r.s.GetCtx().Render(r.viewTemplate, viewData)

		if r.PreRenderHanlder != nil && r.SSRHanlder != nil {
			panic("[inertia]: Can't use pre-render and ssr at the same time, please remove one of the handler (PreRenderHanlder or SSRHanlder)")
		}

		// pre-render and ssr only available on production
		viteDevURL := vite.GetDevelopmentURL()
		if viteDevURL == nil {
			if manifest := vite.GetManifest(); manifest != nil {
				if r.SSRHanlder != nil && r.ssr {
					html := r.s.GetCtx().Response().Body()
					ssr, err := newSSR(html, jsonProps, viewDataStruct)
					if err != nil {
						log.Error(err)
					} else {
						ssr.results = r.SSRHanlder(ssr)

						newHTML, err := ssr.createClientHTML()
						if err != nil {
							log.Error(err)
						}

						if newHTML != nil {
							r.s.GetCtx().Response().SetBody(newHTML)
						}
					}
				} else if r.PreRenderHanlder != nil && r.preRender {
					// file = "assets/app*.js"
					if filePath := manifest.GetCompailedFileNameByResource(rootAppFile); filePath != nil {
						if scriptBuf, err := os.ReadFile(fmt.Sprintf(`./public/build/%s`, *filePath)); err == nil {
							html := string(r.s.GetCtx().Response().Body())

							preRender := newPreRender(&html, filePath, scriptBuf, viewDataStruct, jsonProps)
							preRender.rendered = r.PreRenderHanlder(preRender)

							if preRender.rendered != nil {
								r.s.GetCtx().Response().SetBody(preRender.createClientHTML())
							}
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
