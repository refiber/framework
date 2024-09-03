package refiber

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/refiber/framework/constant"
	"github.com/refiber/framework/internal/validator"
	"github.com/refiber/framework/router"
	"github.com/refiber/framework/support"
)

// TODO: implement init?
// TODO: better error log
// TODO: unit testing

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
				// TODO: use saveTempData (refactor it)
				session, _ := session.Get(c)
				buf, _ := json.Marshal(fiber.Map{
					"type":    support.MessageTypeError,
					"message": message,
				})
				session.Set(string(constant.SessionKeyFlashMessage), buf)
				session.SetExpiry(time.Minute * 1)

				if err := session.Save(); err != nil {
					return err
				}

				return c.RedirectBack("/", 303)
			}

			return c.Status(code).Render("error", fiber.Map{
				"Code":    code,
				"Message": message,
			})
		},
	}

	fiberConfig.Views = newTemplateEngine()

	app := fiber.New(fiberConfig)

	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "header:" + CrsfCookieName,
		CookieName:        CrsfCookieName,
		SingleUseToken:    true,
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		Extractor:         csrf.CsrfFromCookie(CrsfCookieName),
		Session:           session,
		SessionKey:        string(constant.SessionKeyCSRFToken),
	}))

	// TODO: update logger format similar to laravel
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${method} ${latency} ${path}\n",
		TimeFormat: "2 Jan 2006 15:04:05",
	}))

	// avoid placing app.Static before any app.Use, as it will cause the app.Use code or middleware to execute twice.
	// it's fixed on v3: https://github.com/gofiber/fiber/issues/3080
	app.Static("/", "./public")

	/**
	 * Validator config
	 */
	validate := validator.New()

	// TODO: make it support multi lang: https://github.com/go-playground/validator/tree/master/translations
	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)

	/**
	 * Refiber
	 */
	support := support.NewSupport(session, validate, &translator)

	rootRoter := app.Group("")
	router := router.NewRouter(rootRoter, support)

	return app, router, support
}
