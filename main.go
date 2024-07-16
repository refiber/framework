package refiber

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/refiber/framework/constant"
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

func New(c Config) (*fiber.App, router.RouterInterface, support.Refiber) {
	config := configDefault(c)

	/**
	 * Session & crsf
	 */
	session := session.New(session.Config{
		KeyLookup: "cookie:" + SessionName,
		Storage:   config.SessionStorage,
	})

	fiberConfig := fiber.Config{
		AppName: config.AppName,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			}

			debug := os.Getenv("DEBUG")
			if debug == "1" || debug == "true" {
				message = err.Error()
			}

			accepts := c.Accepts("html", "json")
			path := c.Path()
			if accepts == "json" || strings.HasPrefix(path, "/api/") {
				return c.Status(code).JSON(fiber.Map{
					"error": message,
				})
			}

			// handle inertia request
			headers := c.GetReqHeaders()
			if headerXInertia, exist := headers["X-Inertia"]; exist && headerXInertia[0] == "true" {
				// TODO: use saveTempSession (refactor it)
				session, _ := session.Get(c)
				buf, _ := json.Marshal(fiber.Map{
					"type":    support.MessageTypeError,
					"message": message,
				})
				session.Set(constant.KeyFlashMessage+session.ID(), buf)
				session.SetExpiry(time.Minute * 1)

				if err := session.Save(); err != nil {
					return err
				}

				return c.RedirectBack("/", 303)
			}

			return c.Status(code).Render("error", fiber.Map{
				"Code":    code,
				"Message": message,
			}, "error")
		},
	}

	fiberConfig.Views = newTemplateEngine()

	app := fiber.New(fiberConfig)

	app.Static("/", "./public")

	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "header:" + CrsfCookieName,
		CookieName:        CrsfCookieName,
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
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
