package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/refiber/framework/support"
)

type RouterInterface interface {
	Group(path string, middlewares ...Hanlder) *route
	Head(path string, controller Controller, middlewares ...Hanlder) *route
	Get(path string, controller Controller, middlewares ...Hanlder) *route
	Post(path string, controller Controller, middlewares ...Hanlder) *route
	Put(path string, controller Controller, middlewares ...Hanlder) *route
	Delete(path string, controller Controller, middlewares ...Hanlder) *route
	Patch(path string, controller Controller, middlewares ...Hanlder) *route
	CRUD(path string, handler CrudHandler, middlewares ...Hanlder)
}

func NewRouter(rootRoute fiber.Router, support support.Refiber) *route {
	return &route{router: rootRoute, support: support}
}

type route struct {
	router  fiber.Router
	support support.Refiber
}

type Controller = func(support.Refiber, *fiber.Ctx) error

type Hanlder = func(*fiber.Ctx) error

func (r *route) Group(path string, middlewares ...Hanlder) *route {
	return NewRouter(r.router.Group(path, middlewares...), r.support)
}

func (r *route) Get(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Get(path, handlers...)
	return r
}

func (r *route) Head(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Head(path, handlers...)
	return r
}

func (r *route) Post(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Post(path, handlers...)
	return r
}

func (r *route) Put(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Put(path, handlers...)
	return r
}

func (r *route) Delete(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Delete(path, handlers...)
	return r
}

func (r *route) Patch(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support, c)
	})

	r.router.Put(path, handlers...)
	return r
}
