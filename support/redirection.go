package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/constant"
)

func (s *support) Redirect() *redirect {
	return &redirect{s}
}

type redirect struct {
	support *support
}

func (r *redirect) Back() *redirectOptions {
	return &redirectOptions{support: r.support}
}

// TODO: add to external url (https://inertiajs.com/redirects)

func (r *redirect) To(location string) *redirectOptions {
	return &redirectOptions{location: &location, support: r.support}
}

type MessageType string

const (
	MessageTypeInfo    MessageType = "info"
	MessageTypeError   MessageType = "error"
	MessageTypeWarning MessageType = "warning"
	MessageTypeSuccess MessageType = "success"
)

type redirectOptions struct {
	location *string
	support  *support
}

func (ro *redirectOptions) WithMessage(messageType MessageType, message string) *redirectOptions {
	m := fiber.Map{"type": string(messageType), "message": message}

	if err := saveTempData(ro.support, constant.SessionKeyFlashMessage, &m); err != nil {
		log.Errorw("refiber.support.redirection.WithMessage: failed to save session")
	}

	return ro
}

func (ro *redirectOptions) Now() error {
	if ro.location == nil {
		return ro.support.GetCtx().RedirectBack("/", 303)
	}

	return ro.support.GetCtx().Redirect(*ro.location, 303)
}
