package inertia

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/support"
	"github.com/refiber/framework/util"
	"github.com/refiber/framework/vite"
)

type InertiaInterface interface {
	SetViewTemplate(view string)
	Render(*fiber.Ctx) *render
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
		support:                  c.App,
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
	support                  support.Refiber
	viewTemplate             string
	PreRenderHanlder         PreRenderHanlder
	SSRHanlder               SSRHanlder
	EnablePreRenderByDefault bool
	EnableSSRByDefault       bool
}

func (i *inertia) Render(ctx *fiber.Ctx) *render {
	r := &render{}
	r.inertia = i
	r.ctx = ctx
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
	ctx       *fiber.Ctx
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

	if data := r.support.SharedData(r.ctx).GetTemp(); data != nil {
		sharedProps = *data
	}

	data := fiber.Map{}
	data["url"] = r.ctx.OriginalURL()
	v := util.GetMD5Hash("./public/build/manifest.json")
	data["version"] = v
	data["component"] = page
	data["props"] = util.OverrideFiberMaps(&sharedProps, props)

	headers := r.ctx.GetReqHeaders()

	headerXInertia, exist := headers["X-Inertia"]
	headerXInertiaVersion, exist2 := headers["X-Inertia-Version"]
	renderViewTemplate := !exist || !exist2 || len(headerXInertia) > 0 && headerXInertia[0] != "true" || len(headerXInertiaVersion) > 0 && headerXInertiaVersion[0] != v

	if renderViewTemplate {
		jsonProps, _ := json.Marshal(data)
		viewData := createViewData(&jsonProps, r.viewData)

		err := r.ctx.Render(r.viewTemplate, viewData)

		if r.PreRenderHanlder != nil && r.SSRHanlder != nil {
			panic("[inertia]: Can't use pre-render and ssr at the same time, please remove one of the handler (PreRenderHanlder or SSRHanlder)")
		}

		// pre-render and ssr only available on production
		viteDevURL := vite.GetDevelopmentURL()
		if viteDevURL == nil && viewData != nil {
			propsDivTag := (*viewData)["props"].(string)
			if manifest := vite.GetManifest(); manifest != nil {
				if r.SSRHanlder != nil && r.ssr {
					html := r.ctx.Response().Body()
					ssr, err := newSSR(html, jsonProps, &propsDivTag)
					if err != nil {
						log.Error(err)
					} else {
						ssr.results = r.SSRHanlder(ssr)

						newHTML, err := ssr.createClientHTML()
						if err != nil {
							log.Error(err)
						}

						if newHTML != nil {
							r.ctx.Response().SetBody(newHTML)
						}
					}
				} else if r.PreRenderHanlder != nil && r.preRender {
					// file = "assets/app*.js"
					if filePath := manifest.GetCompailedFileNameByResource(rootAppFile); filePath != nil {
						if scriptBuf, err := os.ReadFile(fmt.Sprintf(`./public/build/%s`, *filePath)); err == nil {
							html := string(r.ctx.Response().Body())

							preRender := newPreRender(&html, filePath, scriptBuf, &propsDivTag, jsonProps)
							preRender.rendered = r.PreRenderHanlder(preRender)

							if preRender.rendered != nil {
								r.ctx.Response().SetBody(preRender.createClientHTML())
							}
						}
					}
				}
			}
		}

		return err
	}

	r.ctx.Response().Header.Set("X-Inertia", "true")

	return r.ctx.Status(fiber.StatusOK).JSON(data)
}
