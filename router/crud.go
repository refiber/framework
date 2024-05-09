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

type Crud struct {
	Identifier         string
	Controller         CRUD
	Only               *[]RouteType
	Except             *[]RouteType
	middlewareToRoutes []*crudMiddlewareToRoutes
	availableRoutes    map[RouteType]bool
}

func (c *Crud) fillAvailableRoutes() {
	c.availableRoutes = map[RouteType]bool{
		RouteTypeIndex:   true,
		RouteTypeCreate:  true,
		RouteTypeStore:   true,
		RouteTypeShow:    true,
		RouteTypeEdit:    true,
		RouteTypeUpdate:  true,
		RouteTypeDestroy: true,
	}

	if c.Only == nil && c.Except == nil {
		return
	}

	if c.Only != nil {
		for r := range c.availableRoutes {
			var available bool
			for _, o := range *c.Only {
				if r == o {
					available = true
					continue
				}
			}

			c.availableRoutes[r] = available
		}
	}

	if c.Except != nil {
		for r := range c.availableRoutes {
			available := true
			for _, e := range *c.Except {
				if r == e {
					available = false
					continue
				}
			}

			c.availableRoutes[r] = available
		}
	}
}

func (c *Crud) routeUses(routeType ...RouteType) bool {
	if c.Only == nil && c.Except == nil {
		return true
	}

	for _, r := range routeType {
		if c.availableRoutes[r] {
			return true
		}
	}

	return false
}

func (c *Crud) getMiddlewareForRoute(routeType RouteType) []Hanlder {
	var middlewares = []Hanlder{}

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

type crudMiddlewareToRoutes struct {
	middleware Hanlder
	routeTypes *[]RouteType
}

type CrudHandler = func(crud *Crud)

func (r *route) CRUD(path string, handler CrudHandler, middlewares ...Hanlder) {
	crud := Crud{Identifier: "id"}
	handler(&crud)

	if crud.Controller == nil {
		panic(fmt.Sprintf("[route: %s]: controller is not implemented.", path))
	}

	if crud.Except != nil && crud.Only != nil {
		panic(fmt.Sprintf("[route: %s]: can't use crud.Except and crud.Only in the same time, please choose one.", path))
	}

	crud.fillAvailableRoutes()

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
