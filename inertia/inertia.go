package inertia

import (
	"encoding/json"
	"refiber/support"
	"refiber/utils"

	"github.com/gofiber/fiber/v2"
)

type InertiaInterface interface {
	SetViewTemplate(view string)
	Render() *render
}

func New(s support.Refiber) *inertia {
	return &inertia{s: s, viewTemplate: "app"}
}

type inertia struct {
	s            support.Refiber
	viewTemplate string
}

func (i *inertia) Render() *render {
	r := &render{}
	r.inertia = i
	r.viewTemplate = i.viewTemplate

	return r
}

func (i *inertia) SetViewTemplate(view string) {
	i.viewTemplate = view
}

type render struct {
	*inertia
	viewData *fiber.Map
}

func (r *render) SetViewData(data *fiber.Map) *render {
	r.viewData = data
	return r
}

func (r *render) Page(page string, props *fiber.Map) error {
	sharedProps := fiber.Map{}

	if session, err := r.s.GetSession().Get(r.s.GetCtx()); err == nil {
		sharedProps = *support.GetSharedMap(session)
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

		return r.s.GetCtx().Render(r.viewTemplate, viewData)
	}

	r.s.GetCtx().Response().Header.Set("X-Inertia", "true")

	return r.s.GetCtx().Status(fiber.StatusOK).JSON(data)
}
