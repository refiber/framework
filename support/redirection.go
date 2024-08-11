package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/constant"
)

func (s *support) Redirect(ctx *fiber.Ctx) *redirect {
	return &redirect{s, ctx, s.SharedData(ctx)}
}

type redirect struct {
	support    *support
	ctx        *fiber.Ctx
	sharedData *sharedData
}

func (r *redirect) Back() *redirectOptions {
	return &redirectOptions{redirect: r}
}

// TODO: add to external url (https://inertiajs.com/redirects)

func (r *redirect) To(location string) *redirectOptions {
	return &redirectOptions{redirect: r, location: &location}
}

type MessageType string

const (
	MessageTypeInfo    MessageType = "info"
	MessageTypeError   MessageType = "error"
	MessageTypeWarning MessageType = "warning"
	MessageTypeSuccess MessageType = "success"
)

type redirectOptions struct {
	*redirect
	location *string
}

func (ro *redirectOptions) WithMessage(messageType MessageType, message string) *redirectOptions {
	m := fiber.Map{"type": string(messageType), "message": message}

	if err := ro.sharedData.saveTempData(constant.SessionKeyFlashMessage, &m); err != nil {
		log.Errorw("refiber.support.redirection.WithMessage: failed to save session")
	}

	return ro
}

func (ro *redirectOptions) Now() error {
	if ro.location == nil {
		return ro.ctx.RedirectBack("/", 303)
	}

	return ro.ctx.Redirect(*ro.location, 303)
}
