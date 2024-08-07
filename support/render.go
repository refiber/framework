package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/refiber/framework/utils"
)

func (s *support) Render() *render {
	return &render{s.GetCtx(), s.session}
}

type render struct {
	c       *fiber.Ctx
	session *session.Store
}

func (r *render) View(view string, data *fiber.Map) error {
	m := make(fiber.Map)
	if data != nil {
		m = *data
	}

	s, err := r.session.Get(r.c)
	if err == nil {
		sharedMap := GetTempData(s)
		m = utils.MergeFiberMaps(sharedMap, &m)
	}

	return r.c.Render(view, m)
}

func (r *render) JSON(data *fiber.Map, status int) error {
	m := fiber.Map{}
	if data != nil {
		m = *data
	}

	s, err := r.session.Get(r.c)
	if err == nil {
		sharedMap := GetTempData(s)
		m = utils.MergeFiberMaps(sharedMap, &m)
	}

	return r.c.Status(status).JSON(m)
}
