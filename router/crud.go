package router

import (
	"fmt"

	"github.com/refiber/framework/support"
)

type CRUD interface {
	Index(support.Refiber) error
	Create(support.Refiber) error
	Store(support.Refiber) error
	Show(support.Refiber) error
	Edit(support.Refiber) error
	Update(support.Refiber) error
	Destroy(support.Refiber) error
}

type RouteType int8

const (
	RouteTypeIndex = iota
	RouteTypeCreate
	RouteTypeStore
	RouteTypeShow
	RouteTypeEdit
	RouteTypeUpdate
	RouteTypeDestroy
)

/**
 * For now, customAuthMiddleware is marked as a comment.
 * Initially, I created customAuthMiddleware to add Auth Middleware by default
 * for all routes except RouteTypeIndex and RouteTypeShow. However, upon completion,
 * I realized it's not ideal because you would end up with double Auth middleware
 * if you add route.CRUD into a route group with Auth middleware or add Auth middleware to route.CRUD("/a", func(...), m.Auth).
 */

type Crud struct {
	Identifier         string
	Controller         CRUD
	Only               *[]RouteType
	middlewareToRoutes []*crudMiddlewareToRoutes
	// customAuthMiddleware *Hanlder
}

func (c *Crud) routeUses(routeType ...RouteType) bool {
	if c.Only == nil {
		return true
	}

	for _, v := range *c.Only {
		for _, r := range routeType {
			if v == r {
				return true
			}
		}
	}

	return false
}

func (c *Crud) getMiddlewareForRoute(routeType RouteType) []Hanlder {
	var middlewares = []Hanlder{}

	// if routeType != RouteTypeIndex && routeType != RouteTypeShow {
	// 	middlewares = append(middlewares, *c.customAuthMiddleware)
	// }

	for _, mr := range c.middlewareToRoutes {
		for _, r := range *mr.routeTypes {
			if r == routeType {
				middlewares = append(middlewares, mr.middleware)
			}
		}
	}

	return middlewares
}

func (c *Crud) AddMidlewareToRoutes(middleware Hanlder, routeTypes ...RouteType) {
	c.middlewareToRoutes = append(c.middlewareToRoutes, &crudMiddlewareToRoutes{middleware, &routeTypes})
}

// func (c *Crud) SetAuthMiddleware(m Hanlder) {
// 	c.customAuthMiddleware = &m
// }

type crudMiddlewareToRoutes struct {
	middleware Hanlder
	routeTypes *[]RouteType
}

type CrudHandler = func(crud *Crud)

func (r *route) CRUD(path string, handler CrudHandler, middlewares ...Hanlder) {
	crud := Crud{Identifier: "id"}
	handler(&crud)

	/**
	 * by default create, edit, store, update, and destory are protected by auth middleware
	 * check getMiddlewareForRoute for the logic
	 */
	// if crud.customAuthMiddleware == nil {
	// 	crud.SetAuthMiddleware(func(c *fiber.Ctx) error {
	// 		var user interface{}
	// 		r.support.GetAuthenticatedUserSession(&user)

	// 		if user == nil {
	// 			return support.AuthLoginPage("/login", r.support)
	// 		}

	// 		return c.Next()
	// 	})
	// }

	route := NewRouter(r.router.Group(path, middlewares...), r.support)

	if crud.routeUses(RouteTypeIndex) {
		route.Get("/", crud.Controller.Index)
	}

	if crud.routeUses(RouteTypeCreate) {
		route.Get("/create", crud.Controller.Create, crud.getMiddlewareForRoute(RouteTypeCreate)...)
	}

	if crud.routeUses(RouteTypeStore) {
		route.Post("/create", crud.Controller.Store, crud.getMiddlewareForRoute(RouteTypeStore)...)
	}

	if crud.routeUses(RouteTypeShow, RouteTypeEdit, RouteTypeUpdate, RouteTypeDestroy) {
		routeWithIdentifier := NewRouter(route.router.Group(fmt.Sprintf("/:%s", crud.Identifier)), r.support)

		if crud.routeUses(RouteTypeShow) {
			routeWithIdentifier.Get("/", crud.Controller.Show, crud.getMiddlewareForRoute(RouteTypeShow)...)
		}

		if crud.routeUses(RouteTypeEdit) {
			routeWithIdentifier.Get("/edit", crud.Controller.Edit, crud.getMiddlewareForRoute(RouteTypeEdit)...)
		}

		if crud.routeUses(RouteTypeUpdate) {
			routeWithIdentifier.Put("/edit", crud.Controller.Update, crud.getMiddlewareForRoute(RouteTypeUpdate)...)
		}

		if crud.routeUses(RouteTypeDestroy) {
			routeWithIdentifier.Delete("/delete", crud.Controller.Destroy, crud.getMiddlewareForRoute(RouteTypeDestroy)...)
		}
	}
}
