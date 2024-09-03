package support

import (
	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/util"
)

func (s *support) Render(ctx *fiber.Ctx) *render {
	return &render{ctx, s, s.SharedData(ctx)}
}

type render struct {
	ctx        *fiber.Ctx
	support    *support
	sharedData *sharedData
}

func (r *render) View(view string, data *fiber.Map) error {
	m := make(fiber.Map)
	if data != nil {
		m = *data
	}

	sharedMap := r.sharedData.GetTemp()
	m = util.OverrideFiberMaps(sharedMap, &m)

	return r.ctx.Render(view, m)
}

func (r *render) JSON(data *fiber.Map, status int) error {
	m := fiber.Map{}
	if data != nil {
		m = *data
	}

	sharedMap := r.sharedData.GetTemp()
	m = util.OverrideFiberMaps(sharedMap, &m)

	return r.ctx.Status(status).JSON(m)
}
