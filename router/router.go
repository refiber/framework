package router

import (
	"refiber/support"

	"github.com/gofiber/fiber/v2"
)

type RouterInterface interface {
	Group(path string, middlewares ...Hanlder) *route
	Head(path string, controller Controller, middlewares ...Hanlder) *route
	Get(path string, controller Controller, middlewares ...Hanlder) *route
	Post(path string, controller Controller, middlewares ...Hanlder) *route
	Put(path string, controller Controller, middlewares ...Hanlder) *route
	Delete(path string, controller Controller, middlewares ...Hanlder) *route
	Patch(path string, controller Controller, middlewares ...Hanlder) *route
}

func NewRouter(rootRoute fiber.Router, support support.Refiber) *route {
	return &route{router: rootRoute, support: support}
}

type route struct {
	router  fiber.Router
	support support.Refiber
}

type Controller = func(support.Refiber) error

type Hanlder = func(*fiber.Ctx) error

func (r *route) Group(path string, middlewares ...Hanlder) *route {
	return NewRouter(r.router.Group(path, middlewares...), r.support)
}

func (r *route) Get(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Get(path, handlers...)
	return r
}

func (r *route) Head(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Head(path, handlers...)
	return r
}

func (r *route) Post(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Post(path, handlers...)
	return r
}

func (r *route) Put(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Put(path, handlers...)
	return r
}

func (r *route) Delete(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Put(path, handlers...)
	return r
}

func (r *route) Patch(path string, controller Controller, middlewares ...Hanlder) *route {
	handlers := middlewares
	handlers = append(handlers, func(c *fiber.Ctx) error {
		return controller(r.support)
	})

	r.router.Put(path, handlers...)
	return r
}
