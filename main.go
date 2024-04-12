package refiber

import (
	"html/template"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/refiber/framework/router"
	"github.com/refiber/framework/support"
)

// TODO: implement init?
// TODO: better error log
// TODO: testing

const (
	CrsfCookieName = "XSRF-TOKEN"
	SessionName    = "session"
)

func New(config fiber.Config) (*fiber.App, router.RouterInterface, support.Refiber) {
	if config.AppName == "" {
		config.AppName = "Refiber"
	}

	if config.Views == nil {
		engine := html.New("./resources/views", ".html")
		engine.AddFunc(
			"raw", func(s string) template.HTML {
				return template.HTML(s)
			},
		)
		config.Views = engine
	}

	app := fiber.New(config)

	app.Static("/", "./public")

	/**
	 * Session & crsf
	 */
	session := session.New(session.Config{
		KeyLookup: "cookie:" + SessionName,
	})

	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "header:" + CrsfCookieName,
		CookieName:        CrsfCookieName,
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			accepts := c.Accepts("html", "json")
			path := c.Path()
			if accepts == "json" || strings.HasPrefix(path, "/api/") {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Forbidden",
				})
			}

			// TODO: check if request form inertia, then use flash message

			return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
				"Title":  "Forbidden",
				"Status": fiber.StatusForbidden,
			}, "error")
		},
		Extractor:         csrf.CsrfFromCookie(CrsfCookieName),
		Session:           session,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}))

	/**
	 * Validator config
	 */
	en := en.New()
	uni := ut.New(en, en)

	validate := validator.New(validator.WithRequiredStructEnabled())

	// TODO: make it customizable
	translator, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)

	/**
	 * Refiber
	 */
	support := support.NewSupport(session, validate, &translator)

	/**
	 * Router
	 */
	app.Use(func(c *fiber.Ctx) error {
		support.Ctx = c
		return c.Next()
	})

	rootRoter := app.Group("")
	router := router.NewRouter(rootRoter, support)

	return app, router, support
}
